package workspace

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/meian/atgo/io"
)

var (
	onceMeta sync.Once
	onceTmpl sync.Once
	onceTask sync.Once
)

func metaDir() string {
	metaDir := filepath.Join(Dir(), ".atgo")
	createOnce(&onceMeta, "metadata directory", metaDir)
	return metaDir
}

func TemplateDir() string {
	tmplDir := filepath.Join(metaDir(), "templates")
	createOnce(&onceTmpl, "template directory", tmplDir)
	return tmplDir
}

func TaskRootDir() string {
	taskDir := filepath.Join(metaDir(), "tasks")
	createOnce(&onceTask, "task directory", taskDir)
	return taskDir
}

func LocalTmpDir() (string, error) {
	dir := filepath.Join(metaDir(), "local-tmp")
	if io.DirExists(dir) {
		if err := os.RemoveAll(dir); err != nil {
			return "", err
		}
	}
	return dir, os.Mkdir(dir, 0755)
}
