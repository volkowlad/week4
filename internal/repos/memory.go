package repos

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"week4/internal/myerr"
)

const (
	statusNew      = "new"
	statusProgress = "in_progress"
	statusDone     = "done"
)

type repMemory struct {
	Task sync.Map
}

func NewMemory() Repository {
	return &repMemory{}
}

func checkStatus(status string) string {
	if status == statusProgress || status == statusDone {
		return status
	}

	return statusNew
}

func (r *repMemory) CreateTask(ctx context.Context, task TaskCreate) (uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return uuid.Nil, errors.Wrap(ctx.Err(), "failed to insert task")
	default:
		if task.Title == "" {
			return uuid.Nil, errors.Wrap(myerr.ErrTitle, "failed to insert task")
		}

		newTask := &Task{
			Id:          uuid.New(),
			Title:       task.Title,
			Description: task.Description,
			Status:      statusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		r.Task.Store(newTask.Id, newTask)

		return newTask.Id, nil
	}
}

func (r *repMemory) GetTask(ctx context.Context, id uuid.UUID) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, errors.Wrap(ctx.Err(), "failed to get task")
	default:
		value, ok := r.Task.Load(id)
		if !ok {
			return Task{}, errors.Wrap(myerr.ErrTaskNotFound, "failed to get task")
		}

		task, ok := value.(*Task)
		if !ok {
			return Task{}, errors.Wrap(myerr.ErrInvalidTaskType, "failed to get task")
		}

		return *task, nil
	}
}

func (r *repMemory) GetAllTasks(ctx context.Context) ([]Task, error) {
	select {
	case <-ctx.Done():
		return []Task{}, errors.Wrap(ctx.Err(), "failed to get all tasks")
	default:
		var tasks []Task

		// Используем Range для итерации по всем элементам в sync.Map
		isEmpty := true
		r.Task.Range(func(key, value interface{}) bool {
			isEmpty = false

			task, ok := value.(*Task)
			if !ok {
				return false
			}
			tasks = append(tasks, *task)

			return true
		})

		if isEmpty {
			return nil, errors.Wrap(myerr.ErrTaskNotFound, "failed to get all tasks")
		}

		return tasks, nil
	}
}

func (r *repMemory) DeleteTask(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed to delete task")
	default:
		if _, ok := r.Task.Load(id); ok {
			r.Task.Delete(id)
			return nil
		}

		return errors.Wrap(myerr.ErrTaskNotFound, "failed to delete task")
	}
}

func (r *repMemory) UpdateTask(ctx context.Context, task UpdateTask, id uuid.UUID) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, errors.Wrap(ctx.Err(), "failed to update task")
	default:
		value, ok := r.Task.Load(id)
		if !ok {
			return Task{}, errors.Wrap(myerr.ErrTaskNotFound, "failed to update task")
		}

		newTask, ok := value.(*Task)
		if !ok {
			return Task{}, errors.Wrap(myerr.ErrInvalidTaskType, "failed to update task")
		}

		if task.Title != "" {
			newTask.Title = task.Title
		}

		if task.Description != "" {
			newTask.Description = task.Description
		}

		status := checkStatus(task.Status)
		newTask.Status = status
		newTask.UpdatedAt = time.Now()

		r.Task.Store(newTask.Id, newTask)

		return *newTask, nil
	}
}
