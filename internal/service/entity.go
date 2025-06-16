package service

import (
	"github.com/volkowlad/week4/internal/repos"
)

// TaskRequest - структура, представляющая тело запроса
type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TaskResponse struct {
	Task repos.Task `json:"task"`
}

type AllTasksResponse struct {
	Tasks []repos.Task `json:"all_tasks"`
}
