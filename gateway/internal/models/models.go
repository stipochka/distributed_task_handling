package models

import (
	"encoding/json"
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Task struct {
	Type    string          `json:"type" validate:"required"`
	Payload json.RawMessage `json:"payload" validate:"required"`
}

type Message struct {
	TaskID    uuid.UUID       `json:"task_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
