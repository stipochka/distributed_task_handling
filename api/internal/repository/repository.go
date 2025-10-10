package repository

import (
	"api/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tasksTable = "task_result"
)

type Repository interface {
	GetTask(ctx context.Context, taskID uuid.UUID) (models.Task, error)
	GetAllTaskByType(ctx context.Context, taskType string) ([]models.Task, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(connStr string) (*PostgresRepository, error) {
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: conn,
	}, nil
}

func (p *PostgresRepository) GetTask(ctx context.Context, taskID uuid.UUID) (models.Task, error) {
	const op = "repository.GetTask"

	task := models.Task{}
	query := fmt.Sprintf("SELECT task_id, type, status, description FROM %s WHERE task_id=$1", tasksTable)
	if err := p.db.QueryRow(ctx, query, taskID).Scan(&task.TaskID, &task.Type, &task.Status, &task.Status); err != nil {
		return task, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (p *PostgresRepository) GetAllTaskByType(ctx context.Context, taskType string) ([]models.Task, error) {
	const op = "repository.GetAllTaskByType"

	tasks := make([]models.Task, 0)

	query := fmt.Sprintf("SELECT task_id, type, status, description FROM %s WHERE type=$1", tasksTable)
	rows, err := p.db.Query(ctx, query, taskType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			taskID          uuid.UUID // check if works
			taskType        string
			taskStatus      string
			taskDescription string
		)

		if err := rows.Scan(&taskID, &taskType, &taskStatus, &taskDescription); err != nil {
			return nil, fmt.Errorf("%s: failed to scan rows: %w", op, err)
		}

		tasks = append(tasks, models.Task{
			TaskID:      taskID,
			Type:        taskType,
			Status:      taskStatus,
			Description: taskDescription,
		})

	}

	return tasks, nil
}
