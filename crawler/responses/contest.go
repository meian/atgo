package responses

import (
	"time"

	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
)

type Contest struct {
	ID         string
	Title      string
	StartAt    time.Time
	Duration   time.Duration
	TargetRate string
}

func (c Contest) ToModel() *models.Contest {
	return &models.Contest{
		ID:         ids.ContestID(c.ID),
		Title:      c.Title,
		StartAt:    c.StartAt,
		Duration:   c.Duration,
		TargetRate: c.TargetRate,
	}
}
