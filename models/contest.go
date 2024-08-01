package models

import (
	"time"

	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/url"
	"gopkg.in/guregu/null.v3"
)

type Contest struct {
	ID         ids.ContestID `gorm:"primaryKey"`
	RatedType  null.String   `gorm:"index"`
	Title      string        `gorm:"not null"`
	StartAt    time.Time     `gorm:"not null"`
	Duration   time.Duration `gorm:"not null"`
	TargetRate string        `gorm:"not null"`
	Loaded     bool          `gorm:"type:integer;not null"`

	ContestTasks []ContestTask `gorm:"constraint:OnDelete:CASCADE"`
}

func init() {
	addMigrateTarget(Contest{})
}

func (c Contest) ContestURL() string {
	return url.ContestURL(c.ID)
}

func (c Contest) ContestTaskURL() string {
	return url.ContestTaskURL(c.ID)
}
