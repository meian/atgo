package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type Task struct {
	ID        string        `gorm:"primaryKey"`
	Title     string        `gorm:"not null"`
	TimeLimit time.Duration `gorm:"type:integer;not null"`
	Memory    int           `gorm:"not null"`
	Score     null.Int
	Loaded    bool `gorm:"type:integer;not null"`
}

func init() {
	addMigrateTarget(Task{})
}
