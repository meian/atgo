package flags

import (
	"log/slog"
)

var (
	Version         string
	DefaultLogLevel string
)

func init() {
	if Version == "" {
		Version = "no version, please set with ldflags"
	}
	if DefaultLogLevel == "" {
		DefaultLogLevel = slog.LevelInfo.String()
	}
}
