package requests

import (
	"fmt"
	"net/url"

	"github.com/meian/atgo/constant"
)

type Submit struct {
	ContestID  string
	TaskID     string
	LanguageID constant.LanguageID
	SourceCode string
	CSRFToken  string
}

func (r Submit) URLValues() url.Values {
	return url.Values{
		"data.TaskScreenName": {r.TaskID},
		"data.LanguageId":     {fmt.Sprint(r.LanguageID)},
		"sourceCode":          {r.SourceCode},
		"csrf_token":          {r.CSRFToken},
	}
}
