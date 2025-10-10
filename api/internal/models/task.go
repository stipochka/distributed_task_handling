package models

import "github.com/google/uuid"

type Task struct {
	TaskID      uuid.UUID
	Type        string
	Status      string
	Description string
}
