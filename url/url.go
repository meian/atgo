package url

import (
	"net/url"
	"strings"
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
)

func LoginURL() string {
	return URL("/login", nil, nil).String()
}

func HomeURL() string {
	return URL(HomePath, nil, nil).String()
}

func ContestArchiveURL() string {
	return URL(ContestArchivePath, nil, nil).String()
}

func ContestURL(contestID string) string {
	pathParams := map[string]string{"contestID": contestID}
	return URL(ContestPath, pathParams, nil).String()
}

func ContestTaskURL(contestID string) string {
	pathParams := map[string]string{"contestID": contestID}
	return URL(ContestTaskPath, pathParams, nil).String()
}

func TaskURL(contestID string, taskID string) string {
	pathParams := map[string]string{"contestID": contestID, "id": taskID}
	return URL(TaskPath, pathParams, nil).String()
}

func MySubmissionURL(contestID string) string {
	pathParams := map[string]string{"contestID": contestID}
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
