package workspace

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/meian/atgo/config"
	"github.com/meian/atgo/io"
)

var (
	onceWork sync.Once
	dir      string
)

func Dir() string {
	onceWork.Do(func() {
		ws := config.Config.Workspace
		if ws == "" || ws == "." || ws == ".." || strings.HasPrefix(ws, "./") || strings.HasPrefix(ws, "../") {
			pwd, err := os.Getwd()
			if err != nil {
				slog.Error(err.Error())
				panic("failed to get current directory")
			}
			if ws == "" {
				ws = pwd
			} else if strings.HasPrefix(ws, "..") {
				pwd = filepath.Dir(pwd)
				ws = strings.Replace(ws, "..", pwd, 1)
			} else {
				ws = strings.Replace(ws, ".", pwd, 1)
			}
		}
		dir = ws
		logger := slog.With("dir", dir)
		if _, err := os.Stat(ws); err != nil {
			logger.Error(err.Error())
			panic("failed to check stat for workspace directory")
		}
		logger.Debug("workspace")
	})
	return dir
}

func createOnce(o *sync.Once, target, dir string) {
	o.Do(func() {
		logger := slog.With("target", target, "dir", dir)
		if io.DirExists(dir) {
			logger.Debug("already exists")
			return
		}
		logger.Debug("not exists yet, creating...")
		if err := os.Mkdir(dir, 0755); err != nil {
			logger.Error(err.Error())
			panic("failed to create")
		}
	})
}
