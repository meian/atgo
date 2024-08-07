package url

import (
	"net/url"
	"strings"

	"github.com/meian/atgo/models/ids"
)

const (
	BaseURL            = "https://atcoder.jp"
	HomePath           = "/home"
	LoginPath          = "/login"
	ContestArchivePath = "/contests/archive"
	ContestPath        = "/contests/{contestID}"
	ContestTaskPath    = "/contests/{contestID}/tasks"
	TaskPath           = "/contests/{contestID}/tasks/{id}"
	SubmitPath         = "/contests/{contestID}/submit"
	MySubmissionPath   = "/contests/{contestID}/submissions/me"
	SettingsPath       = "/settings"
)

func LoginURL() string {
	return URL(LoginPath, nil, nil).String()
}

func HomeURL() string {
	return URL(HomePath, nil, nil).String()
}

func ContestArchiveURL() string {
	return URL(ContestArchivePath, nil, nil).String()
}

func ContestURL(contestID ids.ContestID) string {
	pathParams := map[string]string{"contestID": string(contestID)}
	return URL(ContestPath, pathParams, nil).String()
}

func ContestTaskURL(contestID ids.ContestID) string {
	pathParams := map[string]string{"contestID": string(contestID)}
	return URL(ContestTaskPath, pathParams, nil).String()
}

func TaskURL(contestID ids.ContestID, taskID ids.TaskID) string {
	pathParams := map[string]string{"contestID": string(contestID), "id": string(taskID)}
	return URL(TaskPath, pathParams, nil).String()
}

func SubmitURL(contestID ids.ContestID) string {
	pathParams := map[string]string{"contestID": string(contestID)}
	return URL(SubmitPath, pathParams, nil).String()
}

func MySubmissionURL(contestID ids.ContestID) string {
	pathParams := map[string]string{"contestID": string(contestID)}
	return URL(MySubmissionPath, pathParams, nil).String()
}

func URL(path string, pathParams map[string]string, query Valuer) *url.URL {
	for k, v := range pathParams {
		path = strings.ReplaceAll(path, "{"+k+"}", v)
	}
	url, _ := url.Parse(BaseURL)
	url = url.JoinPath(path)
	if query != nil {
		url.RawQuery = query.URLValues().Encode()
	}
	return url
}
