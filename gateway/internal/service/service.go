package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/internal/models"
	"gateway/internal/producer"
	"time"

	"github.com/google/uuid"
)

//go:generate mockery --name TaskService --output ../mocks --dir . --case=underscore
type TaskService interface {
	SendTask(ctx context.Context, task *models.Task) (uuid.UUID, error)
}

type TaskSender struct {
	producer producer.Producer
}

func NewTaskSender(producer producer.Producer) *TaskSender {
	return &TaskSender{
		producer: producer,
	}
}

func (ts *TaskSender) SendTask(ctx context.Context, task *models.Task) (uuid.UUID, error) {
	const op = "service.SendTask"

	taskTime := time.Now()
	id := uuid.New()

	message, err := json.Marshal(models.Message{
		TaskID:    id,
		Type:      task.Type,
		Payload:   task.Payload,
		CreatedAt: taskTime,
	})
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	err = ts.producer.SendMessage(ctx, message)
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
