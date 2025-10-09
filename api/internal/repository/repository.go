package repository

import (
	"api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	GetTask(taskID uuid.UUID) (models.Task, error)
	GetAllTaskByType(taskName string) ([]models.Task, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(connStr string) (*PostgresRepository, error) {
	return nil, nil
}

func (p *PostgresRepository) GetTask(taskID uuid.UUID) (models.Task, error) {
	return models.Task{}, nil
}

func (p *PostgresRepository) GetAllTaskByType(taskName string) ([]models.Task, error) {
	return nil, nil
}
