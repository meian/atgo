package logs

import (
	"log/slog"
	"strings"

	"github.com/pkg/errors"
)

const (
	LevelNone  = slog.LevelError << 1
	LevelError = slog.LevelError
	LevelWarn  = slog.LevelWarn
	LevelInfo  = slog.LevelInfo
	LevelDebug = slog.LevelDebug
)

var (
	levelLookup map[string]slog.Level
)

func init() {
	levelLookup = map[string]slog.Level{
		"none":  LevelNone,
		"error": LevelError,
		"warn":  LevelWarn,
		"info":  LevelInfo,
		"debug": LevelDebug,
	}
}

func ParseLevel(s string) (slog.Level, error) {
	level, ok := levelLookup[strings.ToLower(s)]
	if !ok {
		return LevelNone, errors.Errorf("invalid log level: %s", s)
	}
	return level, nil
}
