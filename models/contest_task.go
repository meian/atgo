package models

import "github.com/meian/atgo/url"

type ContestTask struct {
	ID        string `gorm:"primaryKey"`
	ContestID string `gorm:"not null"`
	TaskID    string `gorm:"not null"`
	Order     int    `gorm:"not null"`
	Index     string `gorm:"not null"`

	Task Task `gorm:"constraint:OnDelete:CASCADE"`
}

func init() {
	addMigrateTarget(ContestTask{})
}

func (ct ContestTask) TaskURL() string {
	return url.TaskURL(ct.ContestID, ct.TaskID)
}
