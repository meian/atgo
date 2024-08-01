package ids

import "regexp"

const taskSampleIDLabel = "task sample ID"

var (
	taskSampleIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*_[a-zA-Z0-9]+__\d+$`)
)

type TaskSampleID string

func (id TaskSampleID) Validate() error {
	if id == "" {
		return ErrEmptyID
	}
	if err := validateLen(taskSampleIDLabel, id); err != nil {
		return err
	}
	if !taskSampleIDPattern.MatchString(string(id)) {
		return newErrInvalidFormat(taskSampleIDLabel, id)
	}
	return nil
}
