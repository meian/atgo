package models

import (
	"fmt"
	"os"

	"github.com/meian/atgo/io"
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
