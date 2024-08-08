package usecase

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/repo"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
)

type Init struct{}

type InitParam struct {
	Force bool
}

type InitResult struct {
	AlreadInitialized bool
	DBFile            string
}

func (u Init) Run(ctx context.Context, param InitParam) (*InitResult, error) {
	logger := logs.FromContext(ctx)

	dbfile, exists := workspace.DBFile()
	logger = logger.With("dbfile", dbfile)

	if exists {
		if param.Force {
			logger.Info("force initialize")
			if _, err := (Clean{}).Run(ctx); err != nil {
				logger.Error(err.Error())
				return nil, errors.New("failed to clean up")
			}
		} else {
			return &InitResult{
				AlreadInitialized: true,
				DBFile:            dbfile,
			}, nil
		}
	}
	logger.Info("start creating db")
	db := database.New(ctx)
	logger.Info("initialized new db file")

	logger.Info("start creating tables")
	err := db.AutoMigrate(models.MigrateTargets...)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create tables")
	}
	tns := make([]string, len(models.MigrateTargets))
	for i, t := range models.MigrateTargets {
		tns[i] = database.TableName(db, t)
	}
	logger.With("tables", tns).Info("create new tables")

	logger.Info("created rated types")
	ratedTypes := []models.RatedType{{Type: "abc"}, {Type: "arc"}, {Type: "agc"}, {Type: "ahc"}}
	if err := repo.NewRateType(db).CreateBatch(ctx, ratedTypes); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create rated types")
	}
	ss := make([]ids.RatedType, len(ratedTypes))
	for i, rt := range ratedTypes {
		ss[i] = rt.Type
	}
	logger.With("rated types", ss).Info("created rated types")

	return &InitResult{
		AlreadInitialized: false,
		DBFile:            dbfile,
	}, nil
}
