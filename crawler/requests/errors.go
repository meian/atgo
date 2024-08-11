package requests

import (
	"github.com/meian/atgo/constant"
	"github.com/pkg/errors"
)

var (
	ErrReqSourceCode = newRequiredError("Source code")
	ErrReqCSRFToken  = newRequiredError("CSRF token")

	ErrPageGT0 = errors.New("Page must be greater than 0")

	ErrInvalidRatedType  = newInvalidErrorFunc[constant.RatedType]("rated type")
	ErrInvalidCategory   = newInvalidErrorFunc[constant.ContestCategory]("category")
	ErrInvalidLanguageID = newInvalidErrorFunc[constant.LanguageID]("language ID")
)

func newRequiredError(name string) error {
	return errors.Errorf("%s is required", name)
}

func newInvalidErrorFunc[T any](name string) func(value T) error {
	return func(value T) error {
		return errors.Errorf("Invalid %s: %v", name, value)
	}
}
