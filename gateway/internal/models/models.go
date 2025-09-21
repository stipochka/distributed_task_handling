package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Type    string `json:"type" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}

type Message struct {
	TaskID    uuid.UUID `json:"task_id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}
