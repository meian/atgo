package ids

import "slices"

type RatedType string

var (
	ratedTypes = []RatedType{"abc", "arc", "agc", "ahc"}
)

func (t RatedType) Validate() error {
	if t == "" {
		return ErrEmptyID
	}
	if !slices.Contains(ratedTypes, t) {
		return newErrInvalidFormat("rated type", t)
	}
	return nil
}
