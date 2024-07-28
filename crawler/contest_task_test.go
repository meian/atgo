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

func TestContestTask_Do_Request(t *testing.T) {
	tests := []struct {
		name string
		req  *requests.ContestTask
		want requestWant
	}{
		{
			name: "success",
			req: &requests.ContestTask{
				ContestID: "abc123",
			},
			want: requestWant{
				path:   "/contests/abc123/tasks",
				query:  url.Values{},
				body:   url.Values{},
				called: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			client, cFunc := mockRequestClient()
			_, _ = crawler.NewContestTask(client).Do(context.Background(), tt.req)
			method, path, query, body, called := cFunc()
			if !tt.want.called {
				assert.False(called)
				return
			}
			assert.Equal(http.MethodGet, method)
			assert.Equal(tt.want.path, path)
			assert.Equal(tt.want.query, query)
			assert.Equal(tt.want.body, body)
			assert.True(called)
		})
	}
}

func TestContestTask_Do_Response(t *testing.T) {
	m := testHTMLMap(t, "contest_task")

	mem := 1024 * 1024 * 1024

	type want struct {
		err bool
		res *responses.ContestTask
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
				res: &responses.ContestTask{
					ContestID: "abc234",
					Tasks: []responses.ContestTask_Task{
						{
							ID:        "abc234_a",
							Title:     "Weird Function",
							Index:     "A",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_b",
							Title:     "Longest Segment",
							Index:     "B",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_c",
							Title:     "Happy New Year!",
							Index:     "C",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_d",
							Title:     "Prefix K-th Max",
							Index:     "D",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_e",
							Title:     "Arithmetic Number",
							Index:     "E",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_f",
							Title:     "Reordering",
							Index:     "F",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_g",
							Title:     "Divide a Sequence",
							Index:     "G",
							TimeLimit: 2 * time.Second,
							Memory:    mem,
						},
						{
							ID:        "abc234_h",
							Title:     "Enumerate Pairs",
							Index:     "Ex",
							TimeLimit: 4 * time.Second,
							Memory:    mem,
						},
					},
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
			req := &requests.ContestTask{ContestID: "abc234"}
			res, err := crawler.NewContestTask(client).Do(ctx, req)
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
