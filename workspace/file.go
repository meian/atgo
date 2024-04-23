package workspace

import (
	"os"
	"path/filepath"
	"time"
)

func CredentialFile() string {
	return filepath.Join(metaDir(), "credential")
}

func CookieFile() (string, time.Time, bool) {
	file := filepath.Join(metaDir(), "cookie")
	stat, err := os.Stat(file)
	if err != nil {
		return file, time.Time{}, false
	}
	return file, stat.ModTime(), true
}
