package repo

import (
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/models"
	"gorm.io/gorm"
)

type Task struct {
	*repository[models.Task]
}

func NewTask(db *gorm.DB) *Task {
	return NewTaskWithDBConn(database.NewDBConn(db))
}

func NewTaskWithDBConn(dbConn *database.DBConn) *Task {
	return &Task{newRepositoryWithDBConn[models.Task](dbConn)}
}
