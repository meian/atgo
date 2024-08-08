package params

import (
	"context"

	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models/ids"
	"github.com/pkg/errors"
)

type Contest struct {
	*baseParam
	RatedType *ids.RatedType
}

func NewContest() *Contest {
	return &Contest{
		baseParam: newBaseParam(50),
	}
}

func (p *Contest) Validate(ctx context.Context) error {
	if p.RatedType != nil {
		if err := p.RatedType.Validate(); err != nil {
			logs.FromContext(ctx).Error(err.Error())
			return errors.New("invalid rated type")
		}
	}
	if err := p.baseParam.Validate(ctx); err != nil {
		return err
	}
	return nil
}
