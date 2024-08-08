package repo

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/meian/atgo/models/ids"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository[M any, ID ids.ModelID] struct {
	DBConn *database.DBConn
}

func newRepositoryWithDBConn[M any, ID ids.ModelID](dbConn *database.DBConn) *repository[M, ID] {
	return &repository[M, ID]{dbConn}
}

func (r *repository[M, ID]) Find(ctx context.Context, id ID) (*M, error) {
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

func (r *repository[M, ID]) Create(ctx context.Context, model *M) error {
	return r.DBConn.DB().Omit(clause.Associations).Create(model).Error
}

func (r *repository[M, ID]) CreateBatch(ctx context.Context, models []M) error {
	return r.DBConn.DB().Omit(clause.Associations).Create(models).Error
}

func (r *repository[M, ID]) Update(ctx context.Context, model *M) error {
	return r.DBConn.DB().Omit(clause.Associations).Save(model).Error
}

func (r *repository[M, ID]) UpdateWithChilds(ctx context.Context, models *M) error {
	return r.DBConn.DB().Save(models).Error
}

func (r *repository[M, ID]) TableName() string {
	var m M
	return database.TableName(r.DBConn.DB(), m)
}

func (r *repository[M, ID]) Truncate(ctx context.Context) error {
	return r.DBConn.DB().Exec("DELETE FROM " + r.TableName()).Error
}

func (r *repository[M, ID]) Tx(f func(conn *database.DBConn) error) error {
	return r.DBConn.Transaction(f)
}
