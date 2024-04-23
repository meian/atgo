package repo

import (
	"context"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/repo/params"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Contest struct {
	*repository[models.Contest]
}

func NewContest(db *gorm.DB) *Contest {
	return NewContestWithDBConn(database.NewDBConn(db))
}

func NewContestWithDBConn(dbConn *database.DBConn) *Contest {
	return &Contest{newRepositoryWithDBConn[models.Contest](dbConn)}
}

func (r *Contest) Search(ctx context.Context, p *params.Contest) ([]models.Contest, error) {
	if p == nil {
		return nil, errors.New("nil parameter")
	}
	if err := p.Validate(ctx); err != nil {
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("invalid search parameter")
	}
	var contests []models.Contest
	query := r.DBConn.DB()
	if p.RatedType != nil && *p.RatedType != constant.RatedTypeAll.String() {
		query = query.Where("rated_type = ?", p.RatedType)
	}
	query = query.Order("start_at DESC")
	query = p.BuildBaseQuery(query)
	if err := query.Find(&contests).Error; err != nil {
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("failed to search contests")
	}
	return contests, nil
}

func (r *Contest) FindWithTasks(ctx context.Context, id string) (*models.Contest, error) {
	var contest models.Contest
	err := r.DBConn.DB().
		Preload("ContestTasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC")
		}).
		Preload("ContestTasks.Task").
		Where("id = ?", id).
		First(&contest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("failed to find contest")
	}
	return &contest, nil
}

func (r *Contest) FindByIDs(ctx context.Context, ids []string) ([]models.Contest, error) {
	var contests []models.Contest
	if err := r.DBConn.DB().Where("id IN (?)", ids).Find(&contests).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find contests: ids=%v", ids)
	}
	return contests, nil
}

func (r *Contest) Truncate(ctx context.Context) error {
	if err := NewTaskWithDBConn(r.DBConn).Truncate(ctx); err != nil {
		return errors.Wrap(err, "failed to truncate tasks")
	}
	return r.repository.Truncate(ctx)
}
