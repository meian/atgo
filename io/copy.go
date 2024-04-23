package io

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open source file: %s", src)
	}
	defer r.Close()
	w, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create destination file: %s", dst)
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	if err != nil {
		return errors.Wrapf(err, "failed to copy file: %s -> %s", src, dst)
	}
	return nil
}
