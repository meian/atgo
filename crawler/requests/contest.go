package requests

import "github.com/meian/atgo/models/ids"

type Contest struct {
	ContestID ids.ContestID
}

func (r Contest) Validate() error {
	return r.ContestID.Validate()
}
