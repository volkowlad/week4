package myerr

import "errors"

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidTaskType = errors.New("invalid task type")
	ErrTitle           = errors.New("title is required")
)
