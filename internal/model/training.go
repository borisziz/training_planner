package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type (
	Training struct {
		ID          int64          `json:"id"`
		Date        time.Time      `json:"date"`
		Items       []TrainingItem `json:"items"`
		Comment     string         `json:"comment"`
		Result      *string        `json:"result,omitempty"`
		ResultFiles []string       `json:"result_files,omitempty"`
	}

	TrainingItem struct {
		ID         int64              `json:"id"`
		TemplateID int64              `json:"template_id" db:"template_id"`
		Params     TrainingItemParams `json:"params" db:"params"`
	}

	TrainingItemParams struct {
		PowerTemplateParams *PowerTemplateParams `json:"power_template_params,omitempty"`
		CrossTemplateParams *CrossTemplateParams `json:"cross_template_params,omitempty"`
	}

	PowerTemplateParams struct {
		Duration int64 `json:"duration"`
	}

	CrossTemplateParams struct {
		Speed    *string `json:"speed,omitempty"`
		Duration *int64  `json:"duration,omitempty"`
		Distance *int64  `json:"distance,omitempty"`
		MaxPulse *int64  `json:"max_pulse,omitempty"`
	}

	UserTraining struct {
		TrainingID int64 `json:"training_id" db:"training_id"`
		UserID     int64 `json:"user_id" db:"user_id"`
	}

	TrainingItemRelation struct {
		TrainingID     int64 `json:"training_id" db:"training_id"`
		TrainingItemID int64 `json:"training_item_id" db:"training_item_id"`
	}
)

// Scan - десериализует данные из БД в виде JSON в структуру
func (tp *TrainingItemParams) Scan(v interface{}) error {
	switch value := v.(type) {
	case []byte:
		err := json.Unmarshal(value, &tp)
		return err
	case string:
		err := json.Unmarshal([]byte(value), &tp)
		return err
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

// Value - сериализует данные структуры в JSON
func (tp *TrainingItemParams) Value() (driver.Value, error) {
	v, err := json.Marshal(tp)
	if err != nil {
		return nil, err
	}
	return string(v), nil
}
