package requests

import (
	"fmt"
	"net/url"

	"github.com/meian/atgo/constant"
)

type ContestArchive struct {
	Page      int
	RatedType constant.RatedType
	Category  constant.ContestCategory
	Keyword   string
}

func (req ContestArchive) URLValues() url.Values {
	vals := url.Values{}
	if req.Page > 0 {
		vals.Add("page", fmt.Sprint(req.Page))
	}
	vals.Add("ratedType", fmt.Sprintf("%d", req.RatedType))
	vals.Add("category", fmt.Sprintf("%d", req.Category))
	if req.Keyword != "" {
		vals.Add("keyword", req.Keyword)
	}
	return vals
}
