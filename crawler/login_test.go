package crawler_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Do_Request(t *testing.T) {
	req := &requests.Login{
		Username:  "user",
		Password:  "pass",
		CSRFToken: "token",
		Continue:  "ctn",
	}
	want := requestWant{
		path:  "/login",
		query: url.Values{"continue": {"ctn"}},
		body:  url.Values{"username": {"user"}, "password": {"pass"}, "csrf_token": {"token"}},
	}

	assert := assert.New(t)
	client, cFunc := mockRequestClient()
	_, _ = crawler.NewLogin(client).Do(context.Background(), req)
	method, path, query, body := cFunc()
	assert.Equal(http.MethodPost, method)
	assert.Equal(want.path, path)
	assert.Equal(want.query, query)
	assert.Equal(want.body, body)
}

func TestLogin_Do_Response(t *testing.T) {
	m := testHTMLMap(t, "login")

	type want struct {
		err bool
		res *responses.Login
	}
	tests := []struct {
		name    string
		httpRes mockHTTPResponse
		want    want
	}{
		{
			name:    "success",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "success.html"},
			want:    want{res: &responses.Login{LoggedIn: true}},
		},
		{
			name:    "forbidden",
			httpRes: mockHTTPResponse{status: http.StatusForbidden},
			want:    want{err: true},
		},
		{
			name:    "not a html response",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "not-a-html"},
			want:    want{err: false, res: &responses.Login{LoggedIn: false}},
		},
		{
			name:    "timeout",
			httpRes: mockHTTPResponse{timeout: true},
			want:    want{err: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			client := mockResponseClient(tt.httpRes.status, m.Get(tt.httpRes.bodyFile), tt.httpRes.timeout)
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
			assert.Equal(tt.want.res, res)
		})
	}
}
