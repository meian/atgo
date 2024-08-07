package requests

import (
	"fmt"
	"net/url"

	"github.com/meian/atgo/constant"
)

type ContestArchive struct {
	Page      int
	RatedType *constant.RatedType
	Category  *constant.ContestCategory
	Keyword   *string
}

func (r ContestArchive) Validate() error {
	if r.Page <= 0 {
		return ErrPageGT0
	}
	if r.RatedType != nil && !r.RatedType.IsARatedType() {
		return ErrInvalidRatedType(*r.RatedType)
	}
	if r.Category != nil && !r.Category.IsAContestCategory() {
		return ErrInvalidCategory(*r.Category)
	}
	return nil
}

func (r ContestArchive) URLValues() url.Values {
	vals := url.Values{}
	if r.Page > 0 {
		vals.Add("page", fmt.Sprint(r.Page))
	}
	if r.RatedType != nil {
		vals.Add("ratedType", fmt.Sprintf("%d", *r.RatedType))
	}
	if r.Category != nil {
		vals.Add("category", fmt.Sprintf("%d", *r.Category))
	}
	if r.Keyword != nil && len(*r.Keyword) > 0 {
		vals.Add("keyword", *r.Keyword)
	}
	return vals
}
