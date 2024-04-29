package usecase

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/meian/atgo/files"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/tmpl"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
)

type TaskLocalInit struct{}

type TaskLocalInitParam struct {
	TaskID    string
	ContestID string
}

type TaskLocalInitResult struct {
	ContestID string
	TaskID    string
	New       bool
}

func (u TaskLocalInit) Run(ctx context.Context, param TaskLocalInitParam) (*TaskLocalInitResult, error) {
	logger := logs.FromContext(ctx)

	tp := TaskParam{
		TaskID:      param.TaskID,
		ContestID:   param.ContestID,
		ShowSamples: true,
	}
	tres, err := Task{}.Run(ctx, tp)
	if err != nil {
		return nil, err
	}

	contest := tres.Contest
	task := tres.Task
	ct := tres.ContestTask
	logger = logger.With("task ID", task.ID).With("contest ID", contest.ID)

	if err := u.backupCurrentWorkspace(ctx); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to backup current workspace")
	}

	info := models.TaskInfo{
		ContestID: contest.ID,
		TaskID:    task.ID,
	}
	if info.CanRestore() {
		logger.With("task directory", info.TaskDir()).
			Info("restore task caused by already exists")
		if err := info.RestoreFiles(); err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to restore task files")
		}
		return &TaskLocalInitResult{
			ContestID: contest.ID,
			TaskID:    task.ID,
			New:       false,
		}, nil
	}

	if err := u.initTemplate(ctx); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to init template")
	}

	logger.Info("create new task files")
	tmpDir, err := os.MkdirTemp("", "atgo")
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create temp work directory")
	}
	defer os.RemoveAll(tmpDir)
	logger = logger.With("temp directory", tmpDir)

	if err := u.executeTaskTemplate(ctx, "main_go", files.MainFile(tmpDir), task, ct); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create main.go")
	}
	if err := u.executeTaskTemplate(ctx, "main_test_go", files.TestFile(tmpDir), task, ct); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create main_test.go")
	}
	if err := u.executeTaskTemplate(ctx, "go_mod", files.ModFile(tmpDir), task, ct); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create go.mod")
	}

	w, err := os.Create(files.TaskLocalFile(tmpDir))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create task file")
	}
	if err := info.Write(w); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to write task file")

	}
	if err := w.Close(); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to close task file")
	}

	// TODO: usecaseで出力は良くないので抑止で良い？
	// それかログに出力する？
	// やり方を考える
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		logger.Error(err.Error())
		logger.Error(stdout.String())
		logger.Error(stderr.String())
		return nil, errors.New("failed to go mod tidy")
	}

	if err := info.CopyFromTemp(tmpDir); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to copy files from temp directory")
	}

	mainfile := files.MainFile(workspace.Dir())
	logger.Info(mainfile)
	cmd = exec.Command("code", mainfile)
	if err := cmd.Run(); err != nil {
		logger.Error(err.Error())
	}

	return &TaskLocalInitResult{
		ContestID: contest.ID,
		TaskID:    task.ID,
		New:       true,
	}, nil
}

func (u TaskLocalInit) backupCurrentWorkspace(ctx context.Context) error {
	var info models.TaskInfo
	file := files.TaskLocalFile(workspace.Dir())
	logger := logs.FromContext(ctx).With("task file", file)
	if !io.FileExists(file) {
		logger.Debug("not found task file")
		return nil
	}
	if err := info.ReadFile(file); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to read task file")
	}
	if !info.RequiredStore() {
		logger.Debug("task is not restorable")
		return nil
	}
	if err := info.StoreFiles(); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to store task files")
	}
	return nil
}

func (u TaskLocalInit) executeTaskTemplate(ctx context.Context, name string, file string, task models.Task, ct models.ContestTask) error {
	tmplFile := workspace.TaskTemplate(name)
	logger := logs.FromContext(ctx).With("template", tmplFile)
	if !io.FileExists(tmplFile) {
		return errors.Errorf("template is not exists: %s", filepath.Base(tmplFile))
	}
	tt, err := tmpl.TaskTemplate(tmplFile)
	if err != nil {
		logger.Error(err.Error())
		return errors.Errorf("failed to parse template: %s", filepath.Base(tmplFile))
	}
	logger = logger.With("file", file)
	w, err := os.Create(file)
	if err != nil {
		logger.Error(err.Error())
		return errors.Errorf("failed to create file: %s", filepath.Base(file))
	}
	defer w.Close()
	data := map[string]any{
		"Task":        task,
		"ContestTask": ct,
	}
	if err := tt.Execute(w, data); err != nil {
		logger.Error(err.Error())
		return errors.Errorf("failed to create file: %s", filepath.Base(file))
	}
	return nil
}

func (u TaskLocalInit) initTemplate(ctx context.Context) error {
	if err := u.initTemplateFile(ctx, "main_go"); err != nil {
		return err
	}
	if err := u.initTemplateFile(ctx, "main_test_go"); err != nil {
		return err
	}
	if err := u.initTemplateFile(ctx, "go_mod"); err != nil {
		return err
	}
	return nil
}

func (u TaskLocalInit) initTemplateFile(ctx context.Context, name string) error {
	file := workspace.TaskTemplate(name)
	logger := logs.FromContext(ctx).With("file", file)
	if io.FileExists(file) {
		return nil
	}
	w, err := os.Create(file)
	if err != nil {
		logger.Error(err.Error())
		return errors.Errorf("failed to create file: %s", filepath.Base(file))
	}
	defer w.Close()
	r := tmpl.TaskTemplateBinary(name)
	if _, err := io.Copy(w, r); err != nil {
		logger.Error(err.Error())
		return errors.Errorf("failed to create template: %s", filepath.Base(file))
	}
	return nil
}