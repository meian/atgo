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
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/repo"
	"github.com/meian/atgo/usecase/common"
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
	info, mustSave, err := u.resolveTaskInfo(ctx, param)
	if err != nil {
		return nil, err
	}
	logger := logs.FromContext(ctx).With("contestID", info.ContestID).With("taskID", info.TaskID)
	ctx = logs.ContextWith(ctx, logger)

	dbConn := database.NewDBConn(database.FromContext(ctx))
	crepo := repo.NewContestWithDBConn(dbConn)
	ctrepo := repo.NewContestTaskWithDBConn(dbConn)
	trepo := repo.NewTaskWithDBConn(dbConn)

	task, err := trepo.Find(ctx, string(info.TaskID))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find task")
	}
	if task == nil {
		task, err = u.createTask(ctx, string(info.ContestID), string(info.TaskID))
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to create task")
		}
	}
	if !task.Loaded {
		err := u.loadTaskSamples(ctx, string(info.ContestID), task)
		if err != nil {
			return nil, err
		}
	}

	contest, err := crepo.Find(ctx, string(info.ContestID))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest")
	}
	if contest == nil {
		contest, err = (&Contest{}).createContest(ctx, string(info.ContestID))
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to create contest")
		}
	}
	ct, err := ctrepo.FindByIDs(ctx, string(info.ContestID), string(info.TaskID))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest task")
	}
	if ct == nil {
		return nil, errors.New("the specified contest and task are not related")
	}
	if param.ShowSamples {
		task, err = trepo.FindWithSamples(ctx, string(info.TaskID))
	} else {
		task.Samples = nil
	}
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find task with samples")
	}

	if mustSave {
		taskFile, _ := workspace.TaskInfoFile()
		if err := info.WriteFile(taskFile); err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to write task info file")
		}
	}

	return &TaskResult{
		Contest:     *contest,
		ContestTask: *ct,
		Task:        *task,
	}, nil

}

func (u Task) resolveTaskInfo(ctx context.Context, param TaskParam) (*models.TaskInfo, bool, error) {
	if len(param.TaskID) == 0 && len(param.ContestID) > 0 {
		return nil, false, errors.New("task ID is required when contest ID is specified")
	}
	logger := logs.FromContext(ctx).With("contestID", param.ContestID).With("taskID", param.TaskID)
	info, mustSave, err := common.ResolveTaskInfo(ctx, param.ContestID, param.TaskID)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, errors.New("task ID cannot be determined")
	}
	if len(info.TaskID) == 0 {
		return nil, false, errors.New("task ID cannot be determined")
	}
	return info, mustSave, nil
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
	req := requests.Task{TaskID: string(task.ID), ContestID: contestID}
	res, err := crawler.NewTask(client).Do(ctx, req)
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
			ID:     ids.TaskSampleID(fmt.Sprintf("%s_%d", task.ID, index)),
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
