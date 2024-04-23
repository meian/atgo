package workspace

import (
	"path/filepath"
	"sync"
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
