package url

import (
	"net/url"
)

const (
	BaseURL   = "https://atcoder.jp"
	HomePath  = "/home"
	LoginPath = "/login"
)

func LoginURL() string {
	return URL("/login", nil).String()
}

func HomeURL() string {
	return URL(HomePath, nil).String()
}

func URL(path string, query Valuer) *url.URL {
	url, _ := url.Parse(BaseURL)
	url = url.JoinPath(path)
	if query != nil {
		url.RawQuery = query.URLValues().Encode()
	}
	return url
}
