package usecase

import (
	"context"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/repo"
	"github.com/meian/atgo/util"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type ContestLoad struct{}

type ContestLoadParam struct {
	RatedType constant.RatedType
	Page      int
}

type ContestLoadResult struct {
	Created     int
	Updated     int
	CurrentPage int
	TotalPages  int
}

func (u ContestLoad) Run(ctx context.Context, param ContestLoadParam) (*ContestLoadResult, error) {
	logger := logs.FromContext(ctx).With("page", param.Page)
	client := http.ClientFromContext(ctx)

	req := requests.ContestArchive{
		Page:      param.Page,
		RatedType: &param.RatedType,
	}
	res, err := crawler.NewContestArchive(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to list contests")
	}

	cids := res.Contests.IDs()
	if len(cids) == 0 {
		return nil, errors.New("no archived contests.")
	}
	logger = logger.With("totalPages", res.TotalPages)
	crepo := repo.NewContest(database.FromContext(ctx))
	contests, err := crepo.FindByIDs(ctx, cids)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to find contests")
	}
	contestm := util.ToMap(contests, func(c models.Contest) ids.ContestID {
		return c.ID
	})

	ratedType := param.RatedType.String()
	var newContests []models.Contest
	var updContests []models.Contest
	for _, contest := range res.Contests {
		logger := logger.With("contestID", contest.ID)
		c, ok := contestm[contest.ID]
		if !ok {
			logger.Debug("new contest")
			newContests = append(newContests, models.Contest{
				ID:         contest.ID,
				Title:      contest.Title,
				StartAt:    contest.StartAt,
				Duration:   contest.Duration,
				RatedType:  null.StringFrom(ratedType),
				TargetRate: contest.TargetRate,
			})
			continue
		}
		if c.RatedType.Valid {
			logger.Debug("no update contest")
			continue
		}
		logger.Debug("update contest")
		c.RatedType = null.StringFrom(ratedType)
		updContests = append(updContests, c)
	}
	if len(newContests) == 0 && len(updContests) == 0 {
		return &ContestLoadResult{
			CurrentPage: req.Page,
			TotalPages:  res.TotalPages,
		}, nil
	}
	err = crepo.Tx(func(conn *database.DBConn) error {
		if len(newContests) > 0 {
			if err := crepo.CreateBatch(ctx, newContests); err != nil {
				logger.Error(err.Error())
				return errors.New("failed to create contests")
			}
			logger.With("count", len(newContests)).Info("created new contests")
		}
		if len(updContests) > 0 {
			for _, c := range updContests {
				logger := logger.With("contestID", c.ID)
				if err := crepo.UpdateWithChilds(ctx, &c); err != nil {
					logger.Error(err.Error())
					return errors.New("failed to update tasks")
				}
				logger.Info("sync contest to rated type")
			}
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return &ContestLoadResult{
		Created:     len(newContests),
		Updated:     len(updContests),
		CurrentPage: req.Page,
		TotalPages:  res.TotalPages,
	}, nil
}
