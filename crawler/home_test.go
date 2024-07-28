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

func TestHome_Do_Request(t *testing.T) {
	req := &requests.Home{}
	want := requestWant{
		path:  "/home",
		query: url.Values{},
		body:  url.Values{},
	}

	assert := assert.New(t)
	client, cFunc := mockRequestClient()
	_, _ = crawler.NewHome(client).Do(context.Background(), req)
	method, path, query, body := cFunc()
	assert.Equal(http.MethodGet, method)
	assert.Equal(want.path, path)
	assert.Equal(want.query, query)
	assert.Equal(want.body, body)
}

func TestHome_Do_Response(t *testing.T) {
	m := testHTMLMap(t, "home")

	type want struct {
		err bool
		res *responses.Home
	}
	tests := []struct {
		name    string
		httpRes mockHTTPResponse
		want    want
	}{
		{
			name:    "success",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "success.html"},
			want:    want{res: &responses.Home{CSRFToken: "test-csrf-token"}},
		},
		{
			name:    "no document token",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "no-document-token.html"},
			want:    want{res: &responses.Home{CSRFToken: ""}},
		},
		{
			name:    "internal server error",
			httpRes: mockHTTPResponse{status: http.StatusInternalServerError},
			want:    want{err: true},
		},
		{
			name:    "not a html response",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "not-a-html"},
			want:    want{err: true},
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
			req := &requests.Home{}
			res, err := crawler.NewHome(client).Do(ctx, req)
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
