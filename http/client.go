package http

import (
	"net/http"
	"net/http/cookiejar"
)

func NewClient(transport http.RoundTripper, jar *cookiejar.Jar) *http.Client {
	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}
}
