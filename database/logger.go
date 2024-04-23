package database

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

var (
	_ logger.Interface = &logDiscard{}
	_ logger.Interface = &dbLogger{}
)

type logDiscard struct{}

func (l *logDiscard) LogMode(logger.LogLevel) logger.Interface    { return l }
func (*logDiscard) Info(context.Context, string, ...interface{})  {}
func (*logDiscard) Warn(context.Context, string, ...interface{})  {}
func (*logDiscard) Error(context.Context, string, ...interface{}) {}
func (*logDiscard) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}

type dbLogger struct {
	*slog.Logger
	level logger.LogLevel
}

func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	// not affected
	return l
}

func (l *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.Logger.Enabled(ctx, slog.LevelDebug) {
		l.Logger.Info(msg, data...)
	}
}

func (l *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.Logger.Enabled(ctx, slog.LevelWarn) {
		l.Logger.Warn(msg, data...)
	}
}

func (l *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.Logger.Enabled(ctx, slog.LevelError) {
		l.Logger.Error(msg, data...)
	}
}

func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	logger := l.Logger
	if l.Logger.Enabled(ctx, slog.LevelDebug) {
		if err != nil {
			logger = l.Logger.With("error", err)
		}
		sql, rows := fc()
		logger.With("rows", rows).With("ellapsed", time.Since(begin)).Debug(sql)
	}
}
