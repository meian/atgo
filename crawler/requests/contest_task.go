package requests

import "github.com/pkg/errors"

type ContestTask struct {
	ContestID string
}

func (r ContestTask) Validate() error {
	if r.ContestID == "" {
		return errors.New("contest id is required")
	}
	return nil
}
