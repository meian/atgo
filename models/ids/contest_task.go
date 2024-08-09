package ids

import (
	"regexp"
)

const contestTaskIDLabel = "contest task ID"

var (
	contestTaskIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*--[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*_[a-zA-Z0-9]+$`)
)

type ContestTaskID string

func (id ContestTaskID) Validate() error {
	if id == "" {
		return ErrEmptyID
	}
	if err := validateLen(contestTaskIDLabel, id); err != nil {
		return err
	}
	if !contestTaskIDPattern.MatchString(string(id)) {
		return newErrInvalidFormat(contestTaskIDLabel, id)
	}
	return nil
}
