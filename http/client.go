package http

import (
	"net/http"
)

func NewClient(transport http.RoundTripper, jar http.CookieJar) *http.Client {
	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}
}
