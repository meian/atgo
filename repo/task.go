package repo

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/pkg/errors"
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

func (r *Task) FindWithSamples(ctx context.Context, taskID string) (*models.Task, error) {
	var task models.Task
	err := r.DBConn.DB().
		Preload("Samples").
		Where("id = ?", taskID).
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("failed to find task")
	}
	return &task, err
}
