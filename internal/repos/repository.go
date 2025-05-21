package repos

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(ctx context.Context, task TaskCreate) (uuid.UUID, error) // Создание задачи
	GetTask(ctx context.Context, id uuid.UUID) (Task, error)
	GetAllTasks(ctx context.Context, page, limit int) ([]Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
	UpdateTask(ctx context.Context, task UpdateTask, id uuid.UUID) (Task, error)
}
