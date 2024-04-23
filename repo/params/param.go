package params

import (
	"context"

	"github.com/meian/atgo/logs"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type baseParam struct {
	Size int
	Page int
}

func newBaseParam(size int) *baseParam {
	return &baseParam{Size: size, Page: 1}
}

func (p *baseParam) Validate(ctx context.Context) error {
	if p.Size <= 0 {
		logs.FromContext(ctx).With("size", p.Size).Error("below 0")
		return errors.New("invalid size")
	}
	if p.Page <= 0 {
		logs.FromContext(ctx).With("page", p.Page).Error("below 0")
		return errors.New("invalid page")
	}
	return nil
}

func (p *baseParam) Limit() int {
	return p.Size
}

func (p *baseParam) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p *baseParam) BuildBaseQuery(query *gorm.DB) *gorm.DB {
	return query.Limit(p.Limit()).Offset(p.Offset())
}
