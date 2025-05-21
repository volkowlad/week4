package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"week4/internal/config"
	"week4/internal/myerr"
)

// SQL-запрос на вставку задачи
const (
	insertTaskQuery  = `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id;`
	selectTasksQuery = `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1;`
	selectAllTasks   = `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	deleteTask       = `DELETE FROM tasks WHERE id = $1;`
)

type repPostgres struct {
	pool *pgxpool.Pool
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewPostgres(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repPostgres{pool}, nil
}

// CreateTask - вставка новой задачи в таблицу tasks
func (r *repPostgres) CreateTask(ctx context.Context, task TaskCreate) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description).Scan(&id)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to insert task")
	}
	return id, nil
}

func (r *repPostgres) GetTask(ctx context.Context, id uuid.UUID) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, selectTasksQuery, id).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return task, myerr.ErrTaskNotFound
		}

		return task, errors.Wrap(err, "failed to query task")
	}

	return task, nil
}

func (r *repPostgres) GetAllTasks(ctx context.Context) ([]Task, error) {
	var tasks []Task

	rows, err := r.pool.Query(ctx, selectAllTasks)
	if err != nil {
		return tasks, errors.Wrap(err, "failed to query tasks")
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return tasks, errors.Wrap(err, "failed to query tasks")
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return tasks, errors.Wrap(err, "failed to query tasks")
	}

	if len(tasks) == 0 {
		return tasks, myerr.ErrTaskNotFound
	}

	return tasks, nil
}

func (r *repPostgres) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, deleteTask, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete task")
	}

	return nil
}

func (r *repPostgres) UpdateTask(ctx context.Context, task UpdateTask, id uuid.UUID) (Task, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	var newTask Task

	if task.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, task.Title)
		argId++
	}

	if task.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, task.Description)
		argId++
	}

	if task.Status == "in_progress" || task.Status == "done" {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, task.Status)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, id)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return newTask, errors.Wrap(err, "failed to start transaction")
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		tx.Rollback(ctx)
		if errors.Cause(err) == pgx.ErrNoRows {
			return newTask, myerr.ErrTaskNotFound
		}
		return newTask, errors.Wrap(err, "failed to update task")
	}

	err = tx.QueryRow(ctx, selectTasksQuery, id).
		Scan(&newTask.Id, &newTask.Title, &newTask.Description, &newTask.Status, &newTask.CreatedAt, &newTask.UpdatedAt)
	if err != nil {
		tx.Rollback(ctx)
		return newTask, errors.Wrap(err, "failed to query task")
	}

	return newTask, tx.Commit(ctx)
}
