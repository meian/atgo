package io

import (
	"os"

	"github.com/meian/atgo/files"
	"github.com/pkg/errors"
)

func Clear(dir string) error {
	if err := clearFile(files.TaskLocalFile(dir)); err != nil {
		return errors.Wrap(err, "failed to clear task file")
	}
	if err := clearFile(files.MainFile(dir)); err != nil {
		return errors.Wrap(err, "failed to clear main file")
	}
	if err := clearFile(files.TestFile(dir)); err != nil {
		return errors.Wrap(err, "failed to clear test file")
	}
	if err := clearFile(files.ModFile(dir)); err != nil {
		return errors.Wrap(err, "failed to clear mod file")
	}
	if err := clearFile(files.SumFile(dir)); err != nil {
		return errors.Wrap(err, "failed to clear sum file")
	}
	tdir := files.TestDataDir(dir)
	if DirExists(tdir) {
		if err := os.RemoveAll(tdir); err != nil {
			return errors.Wrapf(err, "failed to clear testdata: %s", tdir)
		}
	}
	return nil
}

func clearFile(file string) error {
	if !FileExists(file) {
		return nil
	}
	return errors.Wrapf(os.Remove(file), "failed to remove file: %s", file)
}
