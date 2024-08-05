package requests_test

import (
	"net/url"
	"testing"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/crawler/requests"
	"github.com/stretchr/testify/assert"
)

func TestSubmit_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     requests.Submit
		wantErr bool
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
			wantErr: false,
		},
		{
			name: "no contest id",
			req: requests.Submit{
				ContestID:  "",
				TaskID:     "abc123_c",
				CSRFToken:  "token",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "package main\n\nsome code...",
			},
			wantErr: true,
		},
		{
			name: "no task id",
			req: requests.Submit{
				ContestID:  "abc123",
				TaskID:     "",
				CSRFToken:  "token",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "package main\n\nsome code...",
			},
			wantErr: true,
		},
		{
			name: "no language id",
			req: requests.Submit{
				ContestID:  "abc123",
				TaskID:     "abc123_c",
				CSRFToken:  "token",
				LanguageID: constant.LanguageID(0),
				SourceCode: "package main\n\nsome code...",
			},
			wantErr: true,
		},
		{
			name: "invalid language id",
			req: requests.Submit{
				ContestID:  "abc123",
				TaskID:     "abc123_c",
				CSRFToken:  "token",
				LanguageID: constant.LanguageID(100),
				SourceCode: "package main\n\nsome code...",
			},
			wantErr: true,
		},
		{
			name: "no source code",
			req: requests.Submit{
				ContestID:  "abc123",
				TaskID:     "abc123_c",
				CSRFToken:  "token",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSubmit_URLValues(t *testing.T) {
	tests := []struct {
		name string
		req  requests.Submit
		want url.Values
	}{
		{
			name: "submit1",
			req: requests.Submit{
				ContestID:  "abc123",
				TaskID:     "abc123_c",
				CSRFToken:  "token1",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "package main\n\nsome code1...",
			},
			want: url.Values{
				"data.TaskScreenName": {"abc123_c"},
				"data.LanguageId":     {"5002"},
				"sourceCode":          {"package main\n\nsome code1..."},
				"csrf_token":          {"token1"},
			},
		},
		{
			name: "submit2",
			req: requests.Submit{
				ContestID:  "abc456",
				TaskID:     "abc456_c",
				CSRFToken:  "token2",
				LanguageID: constant.LanguageGo_1_20_6,
				SourceCode: "package main\n\nsome code2...",
			},
			want: url.Values{
				"data.TaskScreenName": {"abc456_c"},
				"data.LanguageId":     {"5002"},
				"sourceCode":          {"package main\n\nsome code2..."},
				"csrf_token":          {"token2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.req.URLValues())
		})
	}
}
