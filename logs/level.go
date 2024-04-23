package logs

import (
	"log/slog"

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
	textLookup  map[slog.Level]string
)

func init() {
	levelLookup = map[string]slog.Level{
		"NONE":  LevelNone,
		"ERROR": LevelError,
		"WARN":  LevelWarn,
		"INFO":  LevelInfo,
		"DEBUG": LevelDebug,
	}
	textLookup = make(map[slog.Level]string, len(levelLookup))
	for k, v := range levelLookup {
		textLookup[v] = k
	}
}

func ParseLevel(s string) (slog.Level, error) {
	level, ok := levelLookup[s]
	if !ok {
		return LevelNone, errors.Errorf("invalid log level: %s", s)
	}
	return level, nil
}
