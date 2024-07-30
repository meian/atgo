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
	"github.com/meian/atgo/util"
	"github.com/stretchr/testify/assert"
)

func TestTask_Do_Request(t *testing.T) {
	tests := []struct {
		name string
		req  requests.Task
		want requestWant
	}{
		{
			name: "success",
			req: requests.Task{
				ContestID: "abc123",
				TaskID:    "abc234_d",
			},
			want: requestWant{
				path:   "/contests/abc123/tasks/abc234_d",
				query:  url.Values{},
				body:   url.Values{},
				called: true,
			},
		},
		{
			name: "request is invalid",
			req: requests.Task{
				ContestID: "abc123",
			},
			want: requestWant{called: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			client, cFunc := mockRequestClient()
			_, _ = crawler.NewTask(client).Do(context.Background(), tt.req)
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

func TestTask_Do_Response(t *testing.T) {
	m := testHTMLMap(t, "task")

	type want struct {
		err bool
		res *responses.Task
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
				res: &responses.Task{
					ID:    "abc234_h",
					Score: util.ToPtr(600),
					Samples: []responses.Task_Sample{
						{
							Input:  "6 5\n2 0\n2 2\n3 4\n0 0\n5 5\n8 3\n",
							Output: "9\n1 2\n1 3\n1 4\n2 3\n2 4\n2 5\n3 4\n3 5\n5 6\n",
						},
						{
							Input:  "2 1414213562\n0 0\n1000000000 1000000000\n",
							Output: "0\n",
						},
						{
							Input:  "10 150\n300 300\n300 400\n300 500\n400 300\n400 400\n400 400\n400 500\n500 300\n500 400\n500 500\n",
							Output: "29\n1 2\n1 4\n1 5\n1 6\n2 3\n2 4\n2 5\n2 6\n2 7\n3 5\n3 6\n3 7\n4 5\n4 6\n4 8\n4 9\n5 6\n5 7\n5 8\n5 9\n5 10\n6 7\n6 8\n6 9\n6 10\n7 9\n7 10\n8 9\n9 10\n",
						},
					},
					CSRFToken: "test-csrf-token",
					LoggedIn:  true,
				},
			},
		},
		{
			name:    "no login",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "no-login.html"},
			want: want{
				res: &responses.Task{
					ID:    "abc234_h",
					Score: util.ToPtr(600),
					Samples: []responses.Task_Sample{
						{
							Input:  "6 5\n2 0\n2 2\n3 4\n0 0\n5 5\n8 3\n",
							Output: "9\n1 2\n1 3\n1 4\n2 3\n2 4\n2 5\n3 4\n3 5\n5 6\n",
						},
						{
							Input:  "2 1414213562\n0 0\n1000000000 1000000000\n",
							Output: "0\n",
						},
						{
							Input:  "10 150\n300 300\n300 400\n300 500\n400 300\n400 400\n400 400\n400 500\n500 300\n500 400\n500 500\n",
							Output: "29\n1 2\n1 4\n1 5\n1 6\n2 3\n2 4\n2 5\n2 6\n2 7\n3 5\n3 6\n3 7\n4 5\n4 6\n4 8\n4 9\n5 6\n5 7\n5 8\n5 9\n5 10\n6 7\n6 8\n6 9\n6 10\n7 9\n7 10\n8 9\n9 10\n",
						},
					},
					CSRFToken: "test-csrf-token",
					LoggedIn:  false,
				},
			},
		},
		{
			name:    "no csrf token",
			httpRes: mockHTTPResponse{status: http.StatusOK, bodyFile: "no-token.html"},
			want: want{
				res: &responses.Task{
					ID:    "abc234_h",
					Score: util.ToPtr(600),
					Samples: []responses.Task_Sample{
						{
							Input:  "6 5\n2 0\n2 2\n3 4\n0 0\n5 5\n8 3\n",
							Output: "9\n1 2\n1 3\n1 4\n2 3\n2 4\n2 5\n3 4\n3 5\n5 6\n",
						},
						{
							Input:  "2 1414213562\n0 0\n1000000000 1000000000\n",
							Output: "0\n",
						},
						{
							Input:  "10 150\n300 300\n300 400\n300 500\n400 300\n400 400\n400 400\n400 500\n500 300\n500 400\n500 500\n",
							Output: "29\n1 2\n1 4\n1 5\n1 6\n2 3\n2 4\n2 5\n2 6\n2 7\n3 5\n3 6\n3 7\n4 5\n4 6\n4 8\n4 9\n5 6\n5 7\n5 8\n5 9\n5 10\n6 7\n6 8\n6 9\n6 10\n7 9\n7 10\n8 9\n9 10\n",
						},
					},
					CSRFToken: "",
					LoggedIn:  true,
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
			req := requests.Task{ContestID: "abc234", TaskID: "abc234_h"}
			res, err := crawler.NewTask(client).Do(ctx, req)
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
