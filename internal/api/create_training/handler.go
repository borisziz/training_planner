package create_training

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"training_planner/internal/model"
)

const (
	maxRequestBodySize = 10 << 20 // 10 MiB
	maxMultipartMemory = 1 << 20  // 1 MiB; larger parts use temporary files.
)

type TrainingService interface {
	CreateTraining(ctx context.Context, userID int64, training *model.Training) (int64, error)
}

type createTrainingRequest struct {
	UserID   int64          `json:"user_id"`
	Training model.Training `json:"training"`
}

type Handler struct {
	trainingCreator TrainingService
}

func New(trainingCreator TrainingService) *Handler {
	return &Handler{trainingCreator: trainingCreator}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	request, err := decodeCreateTrainingRequest(r)
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(w, http.StatusRequestEntityTooLarge, "request body is too large")
			return
		}
		var mediaTypeError *unsupportedMediaTypeError
		if errors.As(err, &mediaTypeError) {
			writeError(w, http.StatusUnsupportedMediaType, mediaTypeError.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateCreateTrainingRequest(request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.trainingCreator.CreateTraining(r.Context(), request.UserID, &request.Training)
	if err != nil {
		log.Printf("create training: %v", err)
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create training: %s", err.Error()))
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func decodeCreateTrainingRequest(r *http.Request) (createTrainingRequest, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return createTrainingRequest{}, fmt.Errorf("invalid Content-Type: %w", err)
	}

	switch mediaType {
	case "application/json":
		var request createTrainingRequest
		if err := decodeJSON(r.Body, &request); err != nil {
			return createTrainingRequest{}, fmt.Errorf("invalid JSON body: %w", err)
		}
		return request, nil
	case "multipart/form-data":
		return decodeMultipartCreateTrainingRequest(r)
	default:
		return createTrainingRequest{}, &unsupportedMediaTypeError{mediaType: mediaType}
	}
}

func decodeMultipartCreateTrainingRequest(r *http.Request) (createTrainingRequest, error) {
	if err := r.ParseMultipartForm(maxMultipartMemory); err != nil {
		return createTrainingRequest{}, fmt.Errorf("invalid multipart body: %w", err)
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
		if len(r.MultipartForm.File) > 0 {
			return createTrainingRequest{}, errors.New("file uploads are not supported")
		}
	}

	if payload := strings.TrimSpace(r.FormValue("payload")); payload != "" {
		var request createTrainingRequest
		if err := decodeJSON(strings.NewReader(payload), &request); err != nil {
			return createTrainingRequest{}, fmt.Errorf("invalid payload field: %w", err)
		}
		return request, nil
	}

	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		return createTrainingRequest{}, errors.New("user_id must be an integer")
	}
	var training model.Training
	if err := decodeJSON(strings.NewReader(r.FormValue("training")), &training); err != nil {
		return createTrainingRequest{}, fmt.Errorf("invalid training field: %w", err)
	}
	return createTrainingRequest{UserID: userID, Training: training}, nil
}

func decodeJSON(reader io.Reader, destination any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("body must contain a single JSON value")
		}
		return err
	}
	return nil
}

func validateCreateTrainingRequest(request createTrainingRequest) error {
	if request.UserID <= 0 {
		return errors.New("user_id must be greater than zero")
	}
	if request.Training.Date.IsZero() {
		return errors.New("training.date is required and must use RFC3339 format")
	}
	for index, item := range request.Training.Items {
		if item.TemplateID <= 0 {
			return fmt.Errorf("training.items[%d].template_id must be greater than zero", index)
		}
	}
	return nil
}

type unsupportedMediaTypeError struct {
	mediaType string
}

func (e *unsupportedMediaTypeError) Error() string {
	if e.mediaType == "" {
		return "Content-Type is required"
	}
	return fmt.Sprintf("unsupported Content-Type %q", e.mediaType)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
