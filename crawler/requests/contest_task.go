package requests

type ContestTask struct {
	ContestID string
}

func (r ContestTask) Validate() error {
	if r.ContestID == "" {
		return ErrReqContestID
	}
	return nil
}
