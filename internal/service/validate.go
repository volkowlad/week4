package service

import "github.com/pkg/errors"

func (u UpdateTaskRequest) updateTaskValidate() error {
	if u.Title == "" && u.Description == "" && u.Status == "" {
		err := errors.New("title or description or status is required")
		return errors.Wrap(err, "not validate request to update task")
	}

	return nil
}
