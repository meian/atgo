package usecase

import (
	"context"

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
)

type Contest struct{}

type ContestParam struct {
	ContestID ids.ContestID
}
type ContestResult struct {
	Contest models.Contest
}

func (u Contest) Run(ctx context.Context, param ContestParam) (*ContestResult, error) {
	logger := logs.FromContext(ctx)
	repo := repo.NewContest(database.FromContext(ctx))

	info, mustSave, err := common.ResolveTaskInfo(ctx, param.ContestID, "")
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("no specified contest ID")
	}
	logger = logger.With("contestID", info.ContestID)
	ctx = logs.ContextWith(ctx, logger)

	contest, err := repo.Find(ctx, info.ContestID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest")
	}
	if contest == nil {
		contest, err = u.createContest(ctx, info.ContestID)
		if err != nil {
			return nil, err
		}
	}
	if !contest.Loaded {
		err = u.loadTasks(ctx, contest)
		if err != nil {
			return nil, err
		}
	}

	contest, err = repo.FindWithTasks(ctx, info.ContestID)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contest with tasks")
	}

	if param.ContestID != info.ContestID {
		taskFile, _ := workspace.TaskInfoFile()
		if err := info.WriteFile(taskFile); err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to write task info file")
		}
	}

	if mustSave {
		taskFile, _ := workspace.TaskInfoFile()
		if err := info.WriteFile(taskFile); err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to write task info file")
		}
	}

	return &ContestResult{
		Contest: *contest,
	}, nil
}

func (u Contest) createContest(ctx context.Context, contestID ids.ContestID) (*models.Contest, error) {
	logger := logs.FromContext(ctx)
	client := http.ClientFromContext(ctx)
	req := requests.Contest{ContestID: contestID}
	res, err := crawler.NewContest(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to get contest")
	}
	contest := res.ToModel()
	if err := repo.NewContest(database.FromContext(ctx)).Create(ctx, contest); err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create contest")
	}
	return contest, nil
}

func (u Contest) loadTasks(ctx context.Context, contest *models.Contest) error {
	logger := logs.FromContext(ctx)
	client := http.ClientFromContext(ctx)
	req := requests.ContestTask{ContestID: contest.ID}
	res, err := crawler.NewContestTask(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to get tasks")
	}
	if len(res.Tasks) == 0 {
		return errors.New("no tasks found")
	}
	contest.ContestTasks = res.ToModel()
	contest.Loaded = true
	if err := repo.NewContest(database.FromContext(ctx)).UpdateWithChilds(ctx, contest); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to update contest with tasks")
	}
	return err
}
