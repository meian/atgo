package crawler_test

import (
	"context"
	"net/http"
	gourl "net/url"
	"testing"
	"time"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/url"
	"github.com/stretchr/testify/assert"
)

func TestSubmit_Do_Request(t *testing.T) {
	tests := []struct {
		name string
		req  requests.Submit
		want requestWant
	}{
		{
			name: "success",
			req: requests.Submit{
				ContestID:  "abc123",
				TaskID:     "abc123_c",
				CSRFToken:  "token",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "package main\n\nsome code...",
			},
			want: requestWant{
				path:  "/contests/abc123/submit",
				query: gourl.Values{},
				body: gourl.Values{
					"csrf_token":          []string{"token"},
					"data.LanguageId":     []string{constant.LanguageGo_1_20_6.StringValue()},
					"data.TaskScreenName": []string{"abc123_c"},
					"sourceCode":          []string{"package main\n\nsome code..."},
				},
				called: true,
			},
		},
		{
			name: "request is invalid",
			req: requests.Submit{
				ContestID: "abc123",
			},
			want: requestWant{called: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			client, cFunc := mockRequestClient()
			_, _ = crawler.NewSubmit(client).Do(context.Background(), tt.req)
			method, path, query, body, called := cFunc()
			if !tt.want.called {
				assert.False(called)
				return
			}
			assert.Equal(http.MethodPost, method)
			assert.Equal(tt.want.path, path)
			assert.Equal(tt.want.query, query)
			assert.Equal(tt.want.body, body)
			assert.True(called)
		})
	}
}

func TestSubmit_Do_Response(t *testing.T) {
	m := testHTMLMap(t, "submit")

	type want struct {
		err bool
		res *responses.Submit
	}
	tests := []struct {
		name    string
		httpRes mockHTTPResponse
		want    want
	}{
		{
			name: "success",
			httpRes: mockHTTPResponse{
				lastRequestURL: url.MySubmissionURL("abc123"),
				status:         http.StatusOK,
				bodyFile:       "success.html",
			},
			want: want{
				res: &responses.Submit{Submitted: true},
			},
		},
		{
			name: "invalid parameters",
			httpRes: mockHTTPResponse{
				lastRequestURL: url.SubmitURL("abc123"),
				status:         http.StatusOK,
				bodyFile:       "invalid-parameters.html",
			},
			want: want{
				res: &responses.Submit{Submitted: false},
			},
		},
		{
			name: "not found",
			httpRes: mockHTTPResponse{
				lastRequestURL: url.SubmitURL("abc123"),
				status:         http.StatusNotFound,
			},
			want: want{err: true},
		},
		{
			name: "not a html response",
			httpRes: mockHTTPResponse{
				lastRequestURL: url.SubmitURL("abc123"),
				status:         http.StatusOK,
				bodyFile:       "not-a-html",
			},
			want: want{err: true},
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
			client := tt.httpRes.NewClient(t, m)
			req := requests.Submit{
				ContestID:  "abc123",
				TaskID:     "abc123_c",
				CSRFToken:  "token",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "package main\n\nsome code...",
			}
			res, err := crawler.NewSubmit(client).Do(ctx, req)
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
