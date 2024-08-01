package models

import (
	"time"

	"github.com/meian/atgo/models/ids"
	"gopkg.in/guregu/null.v3"
)

type Task struct {
	ID        ids.TaskID    `gorm:"primaryKey"`
	Title     string        `gorm:"not null"`
	TimeLimit time.Duration `gorm:"type:integer;not null"`
	Memory    int           `gorm:"not null"`
	Score     null.Int
	Loaded    bool `gorm:"type:integer;not null"`

	Samples []TaskSample `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
}

func init() {
	addMigrateTarget(Task{})
}
