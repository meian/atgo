package workspace

import (
	"path/filepath"
	"sync"
)

var (
	onceMeta sync.Once
)

func metaDir() string {
	metaDir := filepath.Join(Dir(), ".atgo")
	createOnce(&onceMeta, "metadata directory", metaDir)
	return metaDir
}
