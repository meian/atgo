package responses

import (
	"time"

	"github.com/meian/atgo/models"
	"gopkg.in/guregu/null.v3"
)

type ContestArchive struct {
	CurrentPage int
	TotalPages  int
	Contests    ContestArchive_ContestList
}

type ContestArchive_ContestList []ContestArchive_Contest

func (c ContestArchive_ContestList) IDs() []string {
	ids := make([]string, len(c))
	for i, contest := range c {
		ids[i] = contest.ID
	}
	return ids
}

type ContestArchive_Contest struct {
	ID         string
	Title      string
	StartAt    time.Time
	Duration   time.Duration
	TargetRate string
}

func (c ContestArchive_Contest) ToModel(ratedType *string) models.Contest {
	return models.Contest{
		ID:         c.ID,
		RatedType:  null.StringFromPtr(ratedType),
		Title:      c.Title,
		StartAt:    c.StartAt,
		Duration:   c.Duration,
		TargetRate: c.TargetRate,
	}
}
