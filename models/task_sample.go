package models

type TaskSampleType int

const (
	TaskSampleTypeSystem TaskSampleType = iota + 1
	TaskSampleTypeUser
)

type TaskSample struct {
	ID     string         `gorm:"primaryKey"`
	TaskID string         `gorm:"not null"`
	Index  string         `gorm:"not null"`
	Input  string         `gorm:"not null"`
	Output string         `gorm:"not null"`
	Type   TaskSampleType `gorm:"not null"`
}

func init() {
	addMigrateTarget(TaskSample{})
}
