package usecase

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/meian/atgo/auth"
	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/files"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
)

type Submit struct{}

type SubmitParam struct {
	SkipTest bool
}

type SubmitErrorStage string

const (
	SubmitErrorStageBuild SubmitErrorStage = "build"
	SubmitErrorStageTest  SubmitErrorStage = "test"
)

type SubmitResult struct {
	LoggedIn     bool
	Submitted    bool
	ErrorStage   SubmitErrorStage
	ErrorMessage string
}

func (u Submit) Run(ctx context.Context, param SubmitParam) (*SubmitResult, error) {
	info, err := u.readTaskLocal(ctx)
	if err != nil {
		return nil, err
	}
	logger := logs.FromContext(ctx).With("contestID", info.ContestID).With("taskID", info.TaskID)

	logger.Info("build before submit")
	if err := u.build(ctx); err != nil {
		logger.With("error", err.Error()).Error("failed to build")
		return &SubmitResult{
			ErrorStage:   SubmitErrorStageBuild,
			ErrorMessage: err.Error(),
		}, nil
	}

	if param.SkipTest {
		logger.Info("skip test")
	} else {
		logger.Info("test before submit")
		if err := u.test(ctx); err != nil {
			logger.With("error", err.Error()).Error("failed to test")
			return &SubmitResult{
				ErrorStage:   SubmitErrorStageTest,
				ErrorMessage: err.Error(),
			}, nil
		}
	}

	source, err := u.sourceCode(ctx)
	if err != nil {
		return nil, err
	}

	loggedIn, token, err := u.login(ctx, *info)
	if err != nil {
		return nil, err
	}
	if !loggedIn {
		return &SubmitResult{
			LoggedIn: false,
		}, nil
	}

	submitted, err := u.submit(ctx, info, source, token)
	if err != nil {
		return nil, err
	}

	return &SubmitResult{
		LoggedIn:  true,
		Submitted: submitted,
	}, nil
}

func (u Submit) readTaskLocal(ctx context.Context) (*models.TaskInfo, error) {
	file := files.TaskLocalFile(workspace.Dir())
	logger := logs.FromContext(ctx).With("local info file", file)
	if !io.FileExists(file) {
		logger.Error("not found task local file")
		return nil, errors.New("not found task local file")
	}
	var info models.TaskInfo
	if err := info.ReadFile(file); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to read task info file")
	}
	return &info, nil
}

func (u Submit) login(ctx context.Context, info models.TaskInfo) (bool, string, error) {
	logger := logs.FromContext(ctx)
	client := http.ClientFromContext(ctx)
	req := requests.Task{
		ContestID: string(info.ContestID),
		TaskID:    string(info.TaskID),
	}
	res, err := crawler.NewTask(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return false, "", errors.New("failed to get task")
	}
	if res.LoggedIn {
		return true, res.CSRFToken, nil
	}
	cred := workspace.CredentialFile()
	if !io.FileExists(cred) {
		logger.Warn("cannot login because credential file is not found")
		return false, "", nil
	}
	username, password, err := auth.Read(ctx, cred)
	if err != nil {
		logger.Error(err.Error())
		return false, "", errors.New("failed to read credential file")
	}
	lreq := &requests.Login{
		Username:  username,
		Password:  password,
		CSRFToken: res.CSRFToken,
		Continue:  info.TaskURL(),
	}
	lres, err := crawler.NewLogin(client).Do(ctx, lreq)
	if err != nil {
		logger.Error(err.Error())
		return false, "", errors.New("failed to login")
	}
	if !lres.LoggedIn {
		logger.Warn("failed to login")
		return false, "", nil
	}
	res, err = crawler.NewTask(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return false, "", errors.New("failed to get task")
	}
	return res.LoggedIn, res.CSRFToken, nil
}

func (u Submit) build(ctx context.Context) error {
	cmd := exec.Command("go", "build", "-o", "/dev/null", files.MainFile(""))
	cmd.Dir = workspace.Dir()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (u Submit) test(context.Context) error {
	cmd := exec.Command("go", "test")
	cmd.Dir = workspace.Dir()
	// テストの構文エラーはstderrで捕捉、テスト自体の失敗はstdoutで捕捉
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stdout.String()
		if msg == "" {
			msg = stderr.String()
		}
		return errors.New(msg)
	}
	return nil
}

func (u Submit) sourceCode(ctx context.Context) (string, error) {
	file := files.MainFile(workspace.Dir())
	logger := logs.FromContext(ctx).With("file", file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		logger.Error(err.Error())
		return "", errors.New("failed to read source code")
	}
	return string(bytes), nil
}

func (u Submit) submit(ctx context.Context, info *models.TaskInfo, source string, csrfToken string) (bool, error) {
	logger := logs.FromContext(ctx)
	client := http.ClientFromContext(ctx)
	req := requests.Submit{
		ContestID:  string(info.ContestID),
		TaskID:     string(info.TaskID),
		LanguageID: constant.LanguageGo_1_20_6,
		SourceCode: source,
		CSRFToken:  csrfToken,
	}
	res, err := crawler.NewSubmit(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return false, errors.New("failed to submit")
	}
	return res.Submitted, nil
}
