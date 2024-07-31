package models

import (
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/url"
)

type ContestTask struct {
	ID        ids.ContestTaskID `gorm:"primaryKey"`
	ContestID ids.ContestID     `gorm:"not null"`
	TaskID    ids.TaskID        `gorm:"not null"`
	Order     int               `gorm:"not null"`
	Index     string            `gorm:"not null"`

	Task Task `gorm:"constraint:OnDelete:CASCADE"`
}

func init() {
	addMigrateTarget(ContestTask{})
}

func (ct ContestTask) TaskURL() string {
	return url.TaskURL(ct.ContestID, ct.TaskID)
}
