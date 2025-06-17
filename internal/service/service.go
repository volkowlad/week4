package service

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/volkowlad/week4/internal/dto"
	"github.com/volkowlad/week4/internal/myerr"
	"github.com/volkowlad/week4/internal/repos"
	"github.com/volkowlad/week4/pkg/validator"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTask(ctx *fiber.Ctx) error
	GetAllTasks(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
}

type service struct {
	repos repos.Repository
	log   *zap.SugaredLogger
}

// NewService - конструктор сервиса
func NewService(repos repos.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repos: repos,
		log:   logger,
	}
}

func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid id")
	}

	// Вставка задачи в БД через репозиторий
	task := repos.TaskCreate{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	}

	taskID, err := s.repos.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Error("Failed to insert task", zap.Error(err))

		if errors.Is(err, myerr.ErrTitle) {
			return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
		}

		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]uuid.UUID{"task_id": taskID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTask(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid id parameter", zap.Error(err))

		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid id parameter")
	}

	var task repos.Task
	task, err = s.repos.GetTask(ctx.Context(), id)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))

		if errors.Is(err, myerr.ErrTaskNotFound) {
			return dto.NotFound(ctx)
		}

		if errors.Is(err, myerr.ErrInvalidTaskType) {
			return dto.WrongType(ctx)
		}

		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   task,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetAllTasks(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)

	if page < 1 {
		page = 1
	}

	limit := ctx.QueryInt("limit", 10)

	if limit < 1 || limit > 100 {
		limit = 10
	}

	var tasks AllTasksResponse
	var err error

	tasks.Tasks, err = s.repos.GetAllTasks(ctx.Context(), page, limit)
	if err != nil {
		s.log.Error("Failed to get all tasks", zap.Error(err))

		if errors.Is(err, myerr.ErrTaskNotFound) {
			return dto.NotFound(ctx)
		}

		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   tasks,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid id parameter", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid id parameter")
	}

	err = s.repos.DeleteTask(ctx.Context(), id)
	if err != nil {
		s.log.Error("Failed to delete task", zap.Error(err))

		if errors.Is(err, myerr.ErrTaskNotFound) {
			return dto.NotFound(ctx)
		}

		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) UpdateTask(ctx *fiber.Ctx) error {
	var req UpdateTaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	err := req.updateTaskValidate()
	if err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, err.Error())
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid id parameter", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid id parameter")
	}

	// Вставка задачи в БД через репозиторий
	task := repos.UpdateTask{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	var newTask TaskResponse

	newTask.Task, err = s.repos.UpdateTask(ctx.Context(), task, id)
	if err != nil {
		s.log.Error("Failed to update task", zap.Error(err))

		if errors.Is(err, myerr.ErrTaskNotFound) {
			return dto.NotFound(ctx)
		}

		if errors.Is(err, myerr.ErrTaskNotFound) {
			return dto.WrongType(ctx)
		}

		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   newTask,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
