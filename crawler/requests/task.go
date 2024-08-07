package requests

type Task struct {
	ContestID string
	TaskID    string
}

func (r Task) Validate() error {
	if r.ContestID == "" {
		return ErrReqContestID
	}
	if r.TaskID == "" {
		return ErrReqTaskID
	}
	return nil
}
