package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/meian/atgo/files"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type taskYAML struct {
	Task TaskInfo `yaml:"atcoder-task"`
}

type TaskInfo struct {
	ContestID string `yaml:"contest-id"`
	TaskID    string `yaml:"task-id"`
}

func (ti TaskInfo) TaskDir() string {
	return filepath.Join(workspace.TaskRootDir(), ti.ContestID, ti.TaskID)
}

func (ti TaskInfo) TaskURL() string {
	return url.TaskURL(ti.ContestID, ti.TaskID)
}

func (ti TaskInfo) RequiredStore() bool {
	return ti.hasRequiredFiles(workspace.Dir())
}

func (ti TaskInfo) CanRestore() bool {
	return ti.hasRequiredFiles(ti.TaskDir())
}

func (ti TaskInfo) hasRequiredFiles(dir string) bool {
	return io.FileExists(files.TaskLocalFile(dir)) &&
		io.FileExists(files.MainFile(dir)) &&
		io.FileExists(files.TestFile(dir)) &&
		io.FileExists(files.ModFile(dir)) &&
		io.FileExists(files.SumFile(dir))
}

func (ti TaskInfo) StoreFiles() error {
	wdir := workspace.Dir()
	if err := ti.WriteFile(files.TaskLocalFile(wdir)); err != nil {
		return errors.Wrap(err, "failed to write task file")
	}
	tdir := ti.TaskDir()
	if err := io.Clear(tdir); err != nil {
		return errors.Wrap(err, "failed to clear task directory")
	}
	return io.Move(wdir, tdir)
}

func (ti TaskInfo) RestoreFiles() error {
	wdir := workspace.Dir()
	tdir := ti.TaskDir()
	if err := io.Clear(wdir); err != nil {
		return errors.Wrap(err, "failed to clear workspace directory")
	}
	return io.Move(tdir, wdir)
}

func (ti TaskInfo) MoveFromTemp(tmpDir string) error {
	return io.Move(tmpDir, workspace.Dir())
}

func (t TaskInfo) Write(w io.Writer) error {
	yamlText := fmt.Sprintf(`atcoder-task:
  contest-id: %s
  task-id: %s
`, t.ContestID, t.TaskID)
	_, err := w.Write([]byte(yamlText))
	return err
}

func (t TaskInfo) WriteFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return errors.Wrapf(err, "failed to create file: %s", file)
	}
	defer f.Close()
	return t.Write(f)
}

func (t *TaskInfo) Read(r io.Reader) error {
	var ty taskYAML
	if err := yaml.NewDecoder(r).Decode(&ty); err != nil {
		return err
	}
	*t = ty.Task
	return nil
}

func (t *TaskInfo) ReadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %s", file)
	}
	defer f.Close()
	return t.Read(f)
}
