package requests

import (
	"net/url"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/models/ids"
)

type Submit struct {
	ContestID  ids.ContestID
	TaskID     ids.TaskID
	LanguageID constant.LanguageID
	SourceCode string
	CSRFToken  string
}

func (r Submit) Validate() error {
	if err := r.ContestID.Validate(); err != nil {
		return err
	}
	if err := r.TaskID.Validate(); err != nil {
		return err
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
		"data.TaskScreenName": {string(r.TaskID)},
		"data.LanguageId":     {r.LanguageID.StringValue()},
		"sourceCode":          {r.SourceCode},
		"csrf_token":          {r.CSRFToken},
	}
}
