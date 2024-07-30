package requests

import "github.com/pkg/errors"

type Task struct {
	ContestID string
	TaskID    string
}

func (r Task) Validate() error {
	if r.ContestID == "" {
		return errors.New("contest id is required")
	}
	if r.TaskID == "" {
		return errors.New("task id is required")
	}
	return nil
}
