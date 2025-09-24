package storage

import (
	"context"
	"fmt"
	"worker/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	resultTable = "task_result"
)

type Storage interface {
	SaveResult(ctx context.Context, result models.Result) error
}

type PostgresStorage struct {
	db *pgxpool.Pool
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(context.Background(), `CREATE TABLE 
		IF NOT EXISTS task_result(
			task_id VARCHAR(32), 
			status VARCHAR(10), 
			description VARCHAR(64)
			PRIMARY KEY task_id;
		);`); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: conn,
	}, nil
}

func (p *PostgresStorage) SaveResult(ctx context.Context, result models.Result) error {
	const op = "storage.TaskResult"

	query := fmt.Sprintf("INSERT INTO %s (task_id, status, description) VALUES ($1, $2, $3);", resultTable)
	_, err := p.db.Exec(ctx, query, result.TaskID, result.Status, result.Description)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
