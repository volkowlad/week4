package repos

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
}

type TaskCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
