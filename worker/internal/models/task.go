package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskID    uuid.UUID       `json:"task_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
