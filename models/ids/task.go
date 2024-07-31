package ids

import "regexp"

const taskIDLabel = "task ID"

var (
	taskIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*_[a-zA-Z0-9]+$`)
)

type TaskID string

func (id TaskID) Validate() error {
	if id == "" {
		return ErrEmptyID
	}
	if err := validateLen(taskIDLabel, id); err != nil {
		return err
	}
	if !taskIDPattern.MatchString(string(id)) {
		return newErrInvalidFormat(taskIDLabel, id)
	}
	return nil
}
