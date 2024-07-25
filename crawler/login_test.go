package crawler_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Do_Request(t *testing.T) {
	req := &requests.Login{
		Username:  "user",
		Password:  "pass",
		CSRFToken: "token",
		Continue:  "ctn",
	}
	want := struct {
		query *url.Values
		body  *url.Values
	}{
		query: &url.Values{"continue": {"ctn"}},
		body:  &url.Values{"username": {"user"}, "password": {"pass"}, "csrf_token": {"token"}},
	}

	assert := assert.New(t)
	client, cFunc := mockRequestClient()
	_, _ = crawler.NewLogin(client).Do(context.Background(), req)
	method, query, body := cFunc()
	assert.Equal(http.MethodPost, method)
	assert.Equal(want.query, query)
	assert.Equal(want.body, body)
}

func TestLogin_Do_Response(t *testing.T) {
	type res struct {
		status   int
		bodyFile string
		timeout  bool
	}
	type want struct {
		err      bool
		loggedIn bool
	}
	tests := []struct {
		name string
		res  res
		want want
	}{
		{
			name: "success",
			res:  res{status: http.StatusOK, bodyFile: "logged-in.html"},
			want: want{loggedIn: true},
		},
		{
			name: "forbidden",
			res:  res{status: http.StatusForbidden},
			want: want{err: true},
		},
		{
			name: "no html",
			res:  res{status: http.StatusOK, bodyFile: "no-html"},
			want: want{err: false, loggedIn: false},
		},
		{
			name: "timeout",
			res:  res{timeout: true},
			want: want{err: true},
		},
	}
	m := testHTMLMap(t, "login")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			client := mockResponseClient(tt.res.status, m.Get(tt.res.bodyFile), tt.res.timeout)
			req := &requests.Login{Username: "user", Password: "pass"}
			res, err := crawler.NewLogin(client).Do(ctx, req)
			if tt.want.err {
				if assert.Error(err) {
					t.Logf("error: %v", err)
				}
				return
			}
			assert.NoError(err)
			if !assert.NotNil(res) {
				return
			}
			assert.Equal(tt.want.loggedIn, res.LoggedIn)
		})
	}
}
