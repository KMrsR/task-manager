package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KMrsR/task-manager/internal/models"
	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &PostgresStorage{conn: conn}, nil
}

// Добавляет новую задачу
func (p *PostgresStorage) AddTask(ctx context.Context, task models.Task) error {
	_, err := p.conn.Exec(ctx,
		`INSERT INTO tasks (id, title, status, deadline) 
		VALUES ($1, $2, $3, $4)`,
		task.ID, task.Title, task.Status, task.Deadline)
	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}
	return nil
}

// Возвращает все задачи
func (p *PostgresStorage) GetTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	sqlStr := `SELECT id,title,status,deadline FROM tasks`
	rows, err := p.conn.Query(ctx, sqlStr)
	if err != nil {
		return nil, fmt.Errorf("failed get tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		var temp time.Time
		err := rows.Scan(&task.ID, &task.Title, &task.Status, &temp)
		task.Deadline = temp.Format(time.RFC1123)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Возвращает задачу по ID
func (p *PostgresStorage) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	sqlStr := `SELECT id,title,status,deadline FROM tasks WHERE id=$1`
	err := p.conn.QueryRow(ctx, sqlStr, id).Scan(task.ID, task.Title, task.Status, task.Deadline)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

// Обновляем задачу по ID
func (p *PostgresStorage) UpdateTask(ctx context.Context, id string, updatedTask models.Task) error {
	sqlStr := `UPDATE tasks SET  title=$1, status=$2, deadline=$3 WHERE id=$4`
	_, err := p.conn.Exec(ctx, sqlStr, updatedTask.Title, updatedTask.Status, updatedTask.Deadline, updatedTask.ID)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

// Удаляем задачу по ID
func (p *PostgresStorage) DeleteTask(ctx context.Context, id string) error {
	sqlStr := `DELETE FROM tasks WHERE id=$1`
	_, err := p.conn.Exec(ctx, sqlStr, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

// закрываем соединение
func (p *PostgresStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.conn.Close(ctx)
}
