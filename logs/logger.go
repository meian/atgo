package logs

import (
	"io"
	"log/slog"

	"github.com/fatih/color"
	"github.com/m-mizutani/clog"
	"github.com/meian/atgo/tmpl"
)

var colorDefault = &clog.ColorMap{
	Level: map[slog.Level]*color.Color{
		slog.LevelInfo:  color.New(color.FgCyan, color.Bold),
		slog.LevelWarn:  color.New(color.FgYellow, color.Bold),
		slog.LevelError: color.New(color.FgRed, color.Bold),
	},
	LevelDefault: color.New(color.BgGreen, color.Bold),
	Time:         color.New(color.FgHiBlack),
	Message:      color.New(color.FgHiBlack, color.Italic),

	AttrKey:   color.New(color.FgHiRed),
	AttrValue: color.New(color.FgHiGreen),
}

func New(out io.Writer, level slog.Level) *slog.Logger {
	if level == LevelNone {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	handler := clog.New(
		clog.WithWriter(out),
		clog.WithLevel(level),
		clog.WithTimeFmt("2006-01-02 15:04:05.999"),
		clog.WithColorMap(colorDefault),
		clog.WithPrinter(clog.IndentPrinter),
		clog.WithSource(true),
		clog.WithTemplate(tmpl.LoggerTemplate()),
	)
	logger := slog.New(handler).WithGroup("params")
	return logger
}
