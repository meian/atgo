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
	"github.com/meian/atgo/timezone"
	"github.com/stretchr/testify/assert"
)

func TestContest_Do_Request(t *testing.T) {
	req := &requests.Contest{
		ContestID: "abc123",
	}
	want := requestWant{
		path:  "/contests/abc123",
		query: url.Values{},
		body:  url.Values{},
	}

	assert := assert.New(t)
	client, cFunc := mockRequestClient()
	_, _ = crawler.NewContest(client).Do(context.Background(), req)
	method, path, query, body := cFunc()
	assert.Equal(http.MethodGet, method)
	assert.Equal(want.path, path)
	assert.Equal(want.query, query)
	assert.Equal(want.body, body)
}

func TestContest_Do_Response(t *testing.T) {
	m := testHTMLMap(t, "contest")

	type want struct {
		err bool
		res *responses.Contest
	}
	tests := []struct {
		name    string
		httpRes mockHTTPResponse
		want    want
	}{
		{
			name:    "success",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "success.html"},
			want: want{
				res: &responses.Contest{
					ID:         "abc234",
					Title:      "AtCoder Beginner Contest 234",
					StartAt:    time.Date(2022, 1, 8, 21, 0, 0, 0, timezone.Tokyo),
					Duration:   1*time.Hour + 40*time.Minute,
					TargetRate: " - 1999",
				},
			},
		},
		{
			name:    "not found",
			httpRes: mockHTTPResponse{status: http.StatusNotFound},
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
			req := &requests.Contest{ContestID: "abc234"}
			res, err := crawler.NewContest(client).Do(ctx, req)
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
