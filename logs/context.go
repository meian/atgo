package logs

import (
	"context"
	"log/slog"
)

type key string

const loggerKey key = "logger"

func FromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey)
	if v == nil {
		return slog.Default()
	}
	return v.(*slog.Logger)
}

func ContextWith(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
