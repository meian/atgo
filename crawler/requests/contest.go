package requests

type Contest struct {
	ContestID string
}

func (r Contest) Validate() error {
	if r.ContestID == "" {
		return ErrReqContestID
	}
	return nil
}
