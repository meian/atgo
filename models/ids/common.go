package ids

import "errors"

type ModelID interface {
	~string
	Validate() error
}

var (
	ErrEmptyID = errors.New("empty id")
)

type errInvalidFormat struct {
	label string
	id    string
}

func (e errInvalidFormat) Error() string {
	return "invalid " + e.label + " format: " + e.id
}

func (e errInvalidFormat) String() string {
	return e.Error()
}

func newErrInvalidFormat[ID ModelID](label string, id ID) error {
	return errInvalidFormat{label: label, id: string(id)}
}

type errTooLong struct {
	label string
	id    string
}

func (e errTooLong) Error() string {
	return e.label + " is too long: " + e.id
}

func (e errTooLong) String() string {
	return e.Error()
}

func validateLen[ID ModelID](label string, id ID) error {
	if len(id) > 64 {
		return errTooLong{label: label, id: string(id)}
	}
	return nil
}
