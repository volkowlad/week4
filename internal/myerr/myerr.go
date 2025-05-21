package myerr

import "github.com/pkg/errors"

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidTaskType = errors.New("invalid task type")
	ErrTitle           = errors.New("title is required")
	ErrRange           = errors.New("page out of range")
)
