package requests

import "github.com/meian/atgo/models/ids"

type ContestTask struct {
	ContestID ids.ContestID
}

func (r ContestTask) Validate() error {
	return r.ContestID.Validate()
}
