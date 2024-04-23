package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(ctx context.Context) *gorm.DB {
	file, _ := workspace.DBFile()
	slogger := logs.FromContext(ctx).With("db", file)
	var glogger logger.Interface
	if slogger.Enabled(ctx, slog.LevelDebug) {
		glogger = &dbLogger{slogger, logger.Info}
	} else {
		glogger = &logDiscard{}
	}
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{
		Logger: glogger,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to open db"))
	}
	if slogger.Enabled(ctx, slog.LevelDebug) {
		db = db.Debug()
	}
	return db
}

func Delete(ctx context.Context) error {
	file, exists := workspace.DBFile()
	if !exists {
		return nil
	}
	if err := os.Remove(file); err != nil {
		logs.FromContext(ctx).Error(err.Error())
		return errors.New("failed to remove db file")
	}
	return nil
}

func NewIfExists(ctx context.Context) (*gorm.DB, bool) {
	_, exists := workspace.DBFile()
	if exists {
		return nil, false
	}
	return New(ctx), true
}
