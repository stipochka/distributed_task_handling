package models

import "github.com/google/uuid"

type Result struct {
	TaskID      uuid.UUID
	Type        string
	Status      string
	Description string
}
