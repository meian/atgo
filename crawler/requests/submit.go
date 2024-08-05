package requests

import (
	"net/url"

	"github.com/meian/atgo/constant"
	"github.com/pkg/errors"
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
		return errors.New("contest id is required")
	}
	if r.TaskID == "" {
		return errors.New("task id is required")
	}
	if !r.LanguageID.Valid() {
		return errors.Errorf("invalid language id: %d", r.LanguageID)
	}
	if r.SourceCode == "" {
		return errors.New("source code is required")
	}
	if r.CSRFToken == "" {
		return errors.New("csrf token is required")
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
