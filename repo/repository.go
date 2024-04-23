package repo

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository[T any] struct {
	DBConn *database.DBConn
}

func newRepositoryWithDBConn[T any](dbConn *database.DBConn) *repository[T] {
	return &repository[T]{dbConn}
}

func (r *repository[T]) Find(ctx context.Context, id string) (*T, error) {
	var m T
	err := r.DBConn.DB().Omit(clause.Associations).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to find record")
	}
	return &m, nil
}

func (r *repository[T]) Create(ctx context.Context, model *T) error {
	return r.DBConn.DB().Omit(clause.Associations).Create(model).Error
}

func (r *repository[T]) CreateBatch(ctx context.Context, models []T) error {
	return r.DBConn.DB().Omit(clause.Associations).Create(models).Error
}

func (r *repository[T]) Update(ctx context.Context, model *T) error {
	return r.DBConn.DB().Omit(clause.Associations).Save(model).Error
}

func (r *repository[T]) Delete(ctx context.Context, model *T) error {
	return r.DBConn.DB().Omit(clause.Associations).Delete(model).Error
}

func (r *repository[T]) TableName() string {
	var m T
	return database.TableName(r.DBConn.DB(), m)
}
