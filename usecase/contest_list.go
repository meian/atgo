package usecase

import (
	"context"

	"github.com/meian/atgo/database"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/repo"
	"github.com/meian/atgo/repo/params"
	"github.com/meian/atgo/util"
	"github.com/pkg/errors"
)

type ContestList struct{}

type ContestListParam struct {
	RatedType string
	Page      int
	Size      int
}

type ContestListResult struct {
	Contests []models.Contest
}

func (u ContestList) Run(ctx context.Context, param ContestListParam) (*ContestListResult, error) {
	p := params.NewContest()
	p.RatedType = util.ToPtr(ids.RatedType(param.RatedType))
	p.Page = param.Page
	p.Size = param.Size
	contests, err := repo.NewContest(database.FromContext(ctx)).Search(ctx, p)
	if err != nil {
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("failed to search contests")
	}

	return &ContestListResult{
		Contests: contests,
	}, nil
}
