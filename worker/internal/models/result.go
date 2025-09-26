package models

import "github.com/google/uuid"

type Result struct {
	TaskID      uuid.UUID
	Status      string
	Description string
}
