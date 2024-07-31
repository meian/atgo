package models

import "github.com/meian/atgo/models/ids"

type TaskSampleType int

const (
	TaskSampleTypeSystem TaskSampleType = iota + 1
	TaskSampleTypeUser
)

type TaskSample struct {
	ID     ids.TaskSampleID `gorm:"primaryKey"`
	TaskID ids.TaskID       `gorm:"not null"`
	Index  string           `gorm:"not null"`
	Input  string           `gorm:"not null"`
	Output string           `gorm:"not null"`
	Type   TaskSampleType   `gorm:"not null"`
}

func init() {
	addMigrateTarget(TaskSample{})
}
