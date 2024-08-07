package requests

import (
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

func (r Submit) Validate() error {
	if r.ContestID == "" {
		return ErrReqContestID
	}
	if r.TaskID == "" {
		return ErrReqTaskID
	}
	if !r.LanguageID.Valid() {
		return ErrInvalidLanguageID(r.LanguageID)
	}
	if r.SourceCode == "" {
		return ErrReqSourceCode
	}
	if r.CSRFToken == "" {
		return ErrReqCSRFToken
	}
	return nil
}

func (r Submit) URLValues() url.Values {
	return url.Values{
		"data.TaskScreenName": {r.TaskID},
		"data.LanguageId":     {r.LanguageID.StringValue()},
		"sourceCode":          {r.SourceCode},
		"csrf_token":          {r.CSRFToken},
	}
}
