package usecase

import (
	"context"
	"fmt"

	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/repo"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type Task struct{}

type TaskParam struct {
	TaskID      string
	ContestID   string
	ShowSamples bool
}

type TaskResult struct {
	Contest     models.Contest
	ContestTask models.ContestTask
	Task        models.Task
}

func (u Task) Run(ctx context.Context, param TaskParam) (*TaskResult, error) {
	logger := logs.FromContext(ctx)
	dbConn := database.NewDBConn(database.FromContext(ctx))
	crepo := repo.NewContestWithDBConn(dbConn)
	ctrepo := repo.NewContestTaskWithDBConn(dbConn)
	trepo := repo.NewTaskWithDBConn(dbConn)

	info, err := u.resolveTaskInfo(ctx, param)
	if err != nil {
		return nil, err
	}
	logger = logger.With("contestID", info.ContestID).With("taskID", info.TaskID)
	ctx = logs.ContextWith(ctx, logger)

	task, err := trepo.Find(ctx, info.TaskID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find task")
	}
	if task == nil {
		task, err = u.createTask(ctx, info.ContestID, info.TaskID)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to create task")
		}
	}
	if !task.Loaded {
		err := u.loadTaskSamples(ctx, info.ContestID, task)
		if err != nil {
			return nil, err
		}
	}

	contest, err := crepo.Find(ctx, info.ContestID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest")
	}
	if contest == nil {
		contest, err = (&Contest{}).createContest(ctx, info.ContestID)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to create contest")
		}
	}
	ct, err := ctrepo.FindByIDs(ctx, info.ContestID, info.TaskID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest task")
	}
	if ct == nil {
		return nil, errors.New("not related contest and task")
	}
	if param.ShowSamples {
		task, err = trepo.FindWithSamples(ctx, info.TaskID)
	} else {
		task.Samples = nil
	}
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find task with samples")
	}

	taskFile, _ := workspace.TaskInfoFile()
	if err := info.WriteFile(taskFile); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to write task info file")
	}

	return &TaskResult{
		Contest:     *contest,
		ContestTask: *ct,
		Task:        *task,
	}, nil

}

func (u Task) resolveTaskInfo(ctx context.Context, param TaskParam) (*models.TaskInfo, error) {
	logger := logs.FromContext(ctx)
	if len(param.TaskID) > 0 {
		logger = logger.With("taskID", param.TaskID)
		if len(param.ContestID) == 0 {
			contestID, err := models.DetectContestID(param.TaskID)
			if err != nil {
				logger.Error(err.Error())
				return nil, errors.New("failed to detect contest ID from task ID")
			}
			param.ContestID = contestID
		}
		return &models.TaskInfo{
			TaskID:    param.TaskID,
			ContestID: param.ContestID,
		}, nil
	}

	info, err := u.readTaskInfo(ctx)
	if err != nil {
		return nil, err
	}
	if len(info.TaskID) == 0 {
		return nil, errors.New("not set task ID in task info file")
	}
	if len(info.ContestID) == 0 {
		return nil, errors.New("not set contest ID in task info file")
	}
	return info, nil
}

func (u Task) readTaskInfo(ctx context.Context) (*models.TaskInfo, error) {
	logger := logs.FromContext(ctx)
	file, ok := workspace.TaskInfoFile()
	logger = logger.With("info file", file)
	if !ok {
		logger.Error("not found task info file")
		return nil, errors.New("not found task info file")
	}
	var info models.TaskInfo
	if err := info.ReadFile(file); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to read task info file")
	}
	return &info, nil
}

func (u Task) createTask(ctx context.Context, contestID string, taskID string) (*models.Task, error) {
	logger := logs.FromContext(ctx)

	ctx = logs.ContextWith(ctx, logger)
	contest, err := repo.NewContest(database.FromContext(ctx)).Find(ctx, contestID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest")
	}
	if contest == nil {
		contest, err = (&Contest{}).createContest(ctx, contestID)
		if err != nil {
			return nil, err
		}
	}
	if !contest.Loaded {
		err = (&Contest{}).loadTasks(ctx, contest)
		if err != nil {
			return nil, err
		}
	}

	task, err := repo.NewTask(database.FromContext(ctx)).Find(ctx, taskID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find task")
	}
	if task == nil {
		return nil, errors.New("not found task")
	}

	return task, nil
}

func (u Task) loadTaskSamples(ctx context.Context, contestID string, task *models.Task) error {
	logger := logs.FromContext(ctx)
	client := http.ClientFromContext(ctx)
	req := &requests.Task{TaskID: task.ID, ContestID: contestID}
	res, err := crawler.NewTaskCrawler(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to get task")
	}
	if res.Score != nil {
		task.Score = null.IntFrom(int64(*res.Score))
	}
	for i, s := range res.Samples {
		index := i + 1
		sample := models.TaskSample{
			ID:     fmt.Sprintf("%s_%d", task.ID, index),
			TaskID: task.ID,
			Index:  fmt.Sprint(i + 1),
			Type:   models.TaskSampleTypeSystem,
			Input:  s.Input,
			Output: s.Output,
		}
		task.Samples = append(task.Samples, sample)
	}
	task.Loaded = true
	if err := repo.NewTask(database.FromContext(ctx)).UpdateWithChilds(ctx, task); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to update task")
	}
	return nil
}
