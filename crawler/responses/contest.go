package responses

import (
	"time"

	"github.com/meian/atgo/models"
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
		ID:         c.ID,
		Title:      c.Title,
		StartAt:    c.StartAt,
		Duration:   c.Duration,
		TargetRate: c.TargetRate,
	}
}
