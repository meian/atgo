package repo

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ContestTask struct {
	*repository[models.ContestTask]
}

func NewContestTask(db *gorm.DB) *ContestTask {
	return NewContestTaskWithDBConn(database.NewDBConn(db))
}

func NewContestTaskWithDBConn(dbConn *database.DBConn) *ContestTask {
	return &ContestTask{newRepositoryWithDBConn[models.ContestTask](dbConn)}
}

func (r *ContestTask) FindByIDs(ctx context.Context, contestID, taskID string) (*models.ContestTask, error) {
	var contestTask models.ContestTask
	err := r.DBConn.DB().
		Preload("Task").
		Where("contest_id = ? AND task_id = ?", contestID, taskID).
		First(&contestTask).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("failed to find contest task")
	}
	return &contestTask, nil
}
