package requests

import "github.com/meian/atgo/models/ids"

type Task struct {
	ContestID ids.ContestID
	TaskID    ids.TaskID
}

func (r Task) Validate() error {
	if err := r.ContestID.Validate(); err != nil {
		return err
	}
	return r.TaskID.Validate()
}
