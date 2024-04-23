package workspace

import (
	"os"
	"path/filepath"
	"time"

	"github.com/meian/atgo/io"
)

func CredentialFile() string {
	return filepath.Join(metaDir(), "credential")
}

func DBFile() (string, bool) {
	file := filepath.Join(metaDir(), "atgo.db")
	return file, io.FileExists(file)
}

func CookieFile() (string, time.Time, bool) {
	file := filepath.Join(metaDir(), "cookie")
	stat, err := os.Stat(file)
	if err != nil {
		return file, time.Time{}, false
	}
	return file, stat.ModTime(), true
}
