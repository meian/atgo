package ids

import "regexp"

const contestIDLabel = "contest ID"

var (
	contestIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*$`)
)

type ContestID string

func (id ContestID) Validate() error {
	if id == "" {
		return ErrEmptyID
	}
	if err := validateLen(contestIDLabel, id); err != nil {
		return err
	}
	if !contestIDPattern.MatchString(string(id)) {
		return newErrInvalidFormat(contestIDLabel, id)
	}
	return nil
}
