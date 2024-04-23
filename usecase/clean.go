package usecase

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
)

type Clean struct{}

type CleanResult struct {
	AlreadyCleaned bool
	Cleaned        bool
}

func (u Clean) Run(ctx context.Context) (*CleanResult, error) {
	logger := logs.FromContext(ctx)

	logger.Info("check database exists")
	dbfile, exists := workspace.DBFile()
	logger = logger.With("dbfile", dbfile)
	if !exists {
		return &CleanResult{AlreadyCleaned: true}, nil
	}

	logger.Info("start cleaning up")
	if err := database.Delete(ctx); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to delete database file")
	}

	return &CleanResult{Cleaned: true}, nil
}
