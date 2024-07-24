package requests_test

import (
	"net/url"
	"testing"

	"github.com/meian/atgo/crawler/requests"
	"github.com/stretchr/testify/assert"
)

func TestLogin_URLValues(t *testing.T) {
	type fields struct {
		Username  string
		Password  string
		CSRFToken string
		Continue  string
	}
	tests := []struct {
		name   string
		fields fields
		want   url.Values
	}{
		{
			name: "user1",
			fields: fields{
				Username:  "user1",
				Password:  "pass1",
				CSRFToken: "token1",
				Continue:  "continue1",
			},
			want: url.Values{
				"username":   {"user1"},
				"password":   {"pass1"},
				"csrf_token": {"token1"},
			},
		},
		{
			name: "user2",
			fields: fields{
				Username:  "user2",
				Password:  "pass2",
				CSRFToken: "token2",
				Continue:  "continue2",
			},
			want: url.Values{
				"username":   {"user2"},
				"password":   {"pass2"},
				"csrf_token": {"token2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := requests.Login{
				Username:  tt.fields.Username,
				Password:  tt.fields.Password,
				CSRFToken: tt.fields.CSRFToken,
				Continue:  tt.fields.Continue,
			}
			assert.Equal(t, tt.want, r.URLValues())
		})
	}
}
