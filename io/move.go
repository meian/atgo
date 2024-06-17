package io

import (
	"os"

	"github.com/meian/atgo/files"
	"github.com/pkg/errors"
)

func Move(sdir, ddir string) error {
	if !DirExists(sdir) {
		return nil
	}
	if !DirExists(ddir) {
		if err := os.MkdirAll(ddir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create destination directory: %s", ddir)
		}
	}
	if err := moveFile(files.TaskLocalFile(sdir), files.TaskLocalFile(ddir)); err != nil {
		return errors.Wrap(err, "failed to move task file")
	}
	if err := moveFile(files.MainFile(sdir), files.MainFile(ddir)); err != nil {
		return errors.Wrap(err, "failed to move main file")
	}
	if err := moveFile(files.TestFile(sdir), files.TestFile(ddir)); err != nil {
		return errors.Wrap(err, "failed to move test file")
	}
	if err := moveFile(files.ModFile(sdir), files.ModFile(ddir)); err != nil {
		return errors.Wrap(err, "failed to move mod file")
	}
	if err := moveFile(files.SumFile(sdir), files.SumFile(ddir)); err != nil {
		return errors.Wrap(err, "failed to move sum file")
	}
	std := files.TestDataDir(sdir)
	if DirExists(std) {
		dtd := files.TestDataDir(ddir)
		if err := os.Rename(std, dtd); err != nil {
			return errors.Wrapf(err, "failed to move testdata: %s -> %s", std, dtd)
		}
	}
	return nil
}

func moveFile(src, dst string) error {
	if !FileExists(src) {
		return nil
	}
	if err := os.Rename(src, dst); err != nil {
		return errors.Wrapf(err, "failed to move file: %s -> %s", src, dst)
	}
	return nil
}
