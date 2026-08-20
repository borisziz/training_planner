package create_training

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"training_planner/internal/model"
)

type healthCheckerStub struct {
	err error
}

func (s healthCheckerStub) PingContext(context.Context) error {
	return s.err
}

type trainingCreatorStub struct {
	id       int64
	err      error
	userID   int64
	training model.Training
	calls    int
}

func (s *trainingCreatorStub) CreateTraining(
	_ context.Context,
	userID int64,
	training model.Training,
) (int64, error) {
	s.calls++
	s.userID = userID
	s.training = training
	return s.id, s.err
}

func TestCreateTrainingJSON(t *testing.T) {
	creator := &trainingCreatorStub{id: 42}
	handler := New(healthCheckerStub{}, creator)
	body := `{
		"user_id": 7,
		"training": {
			"date": "2026-08-19T18:30:00+03:00",
			"comment": "intervals",
			"items": [{
				"template_id": 3,
				"params": {"cross_template_params": {
					"speed": "4:30/km", "duration": 1200, "max_pulse": 170
				}}
			}]
		}
	}`

	request := httptest.NewRequest(http.MethodPost, "/training/create", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if creator.calls != 1 || creator.userID != 7 {
		t.Fatalf("unexpected service call: calls=%d userID=%d", creator.calls, creator.userID)
	}
	if creator.training.Comment != "intervals" || len(creator.training.Items) != 1 {
		t.Fatalf("unexpected training: %#v", creator.training)
	}

	var result map[string]int64
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(result, map[string]int64{"id": 42}) {
		t.Fatalf("unexpected response: %#v", result)
	}
}

func TestCreateTrainingMultipart(t *testing.T) {
	creator := &trainingCreatorStub{id: 99}
	handler := New(healthCheckerStub{}, creator)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("user_id", "8"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("training", `{
		"date":"2026-08-20T07:00:00Z",
		"comment":"easy run",
		"items":[{"template_id":4,"params":{}}]
	}`); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/training/create", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if creator.userID != 8 || creator.training.Comment != "easy run" {
		t.Fatalf("unexpected service arguments: userID=%d training=%#v", creator.userID, creator.training)
	}
}

func TestCreateTrainingRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		status      int
	}{
		{
			name:        "unsupported media type",
			contentType: "text/plain",
			body:        "training",
			status:      http.StatusUnsupportedMediaType,
		},
		{
			name:        "unknown JSON field",
			contentType: "application/json",
			body:        `{"user_id":1,"training":{"date":"2026-08-19T10:00:00Z"},"unknown":true}`,
			status:      http.StatusBadRequest,
		},
		{
			name:        "invalid user",
			contentType: "application/json",
			body:        `{"user_id":0,"training":{"date":"2026-08-19T10:00:00Z"}}`,
			status:      http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			creator := &trainingCreatorStub{id: 1}
			handler := New(healthCheckerStub{}, creator)
			request := httptest.NewRequest(http.MethodPost, "/training/create", bytes.NewBufferString(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.status {
				t.Fatalf("expected status %d, got %d: %s", test.status, response.Code, response.Body.String())
			}
			if creator.calls != 0 {
				t.Fatalf("service must not be called, got %d calls", creator.calls)
			}
		})
	}
}

func TestCreateTrainingHidesServiceError(t *testing.T) {
	creator := &trainingCreatorStub{err: errors.New("database credentials leaked")}
	handler := New(healthCheckerStub{}, creator)
	body := `{"user_id":1,"training":{"date":"2026-08-19T10:00:00Z"}}`
	request := httptest.NewRequest(http.MethodPost, "/training/create", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
	if bytes.Contains(response.Body.Bytes(), []byte("credentials")) {
		t.Fatalf("internal error was exposed: %s", response.Body.String())
	}
}

func TestHealth(t *testing.T) {
	handler := New(healthCheckerStub{}, &trainingCreatorStub{})
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestTrainingDateUsesRFC3339(t *testing.T) {
	var training model.Training
	if err := json.Unmarshal([]byte(`{"date":"2026-08-19T10:00:00Z"}`), &training); err != nil {
		t.Fatal(err)
	}
	if !training.Date.Equal(time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected date: %s", training.Date)
	}
}
