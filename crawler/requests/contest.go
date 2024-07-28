package requests

import "github.com/pkg/errors"

type Contest struct {
	ContestID string
}

func (r Contest) Validate() error {
	if r.ContestID == "" {
		return errors.New("contest id is required")
	}
	return nil
}
