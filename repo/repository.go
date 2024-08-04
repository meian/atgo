package repo

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository[M any] struct {
	DBConn *database.DBConn
}

func newRepositoryWithDBConn[M any](dbConn *database.DBConn) *repository[M] {
	return &repository[M]{dbConn}
}

func (r *repository[M]) Find(ctx context.Context, id string) (*M, error) {
	var m M
	err := r.DBConn.DB().Omit(clause.Associations).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to find record")
	}
	return &m, nil
}

func (r *repository[M]) Create(ctx context.Context, model *M) error {
	return r.DBConn.DB().Omit(clause.Associations).Create(model).Error
}

func (r *repository[M]) CreateBatch(ctx context.Context, models []M) error {
	return r.DBConn.DB().Omit(clause.Associations).Create(models).Error
}

func (r *repository[M]) Update(ctx context.Context, model *M) error {
	return r.DBConn.DB().Omit(clause.Associations).Save(model).Error
}

func (r *repository[M]) UpdateWithChilds(ctx context.Context, models *M) error {
	return r.DBConn.DB().Save(models).Error
}

func (r *repository[M]) TableName() string {
	var m M
	return database.TableName(r.DBConn.DB(), m)
}

func (r *repository[M]) Truncate(ctx context.Context) error {
	return r.DBConn.DB().Exec("DELETE FROM " + r.TableName()).Error
}

func (r *repository[M]) Tx(f func(conn *database.DBConn) error) error {
	return r.DBConn.Transaction(f)
}
