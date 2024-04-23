package params

import (
	"context"
	"slices"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/logs"
	"github.com/pkg/errors"
)

type Contest struct {
	*baseParam
	RatedType *string
}

func NewContest() *Contest {
	return &Contest{
		baseParam: newBaseParam(50),
	}
}

func (p *Contest) Validate(ctx context.Context) error {
	if p.RatedType != nil && !slices.Contains(constant.RatedTypeStrings(), *p.RatedType) {
		logs.FromContext(ctx).With("ratedType", p.RatedType).Error("undefined type")
		return errors.New("invalid rated type")
	}
	if err := p.baseParam.Validate(ctx); err != nil {
		return err
	}
	return nil
}
