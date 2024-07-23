package csrf_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/meian/atgo/csrf"
	"github.com/stretchr/testify/assert"
)

func TestFromCookies(t *testing.T) {

	dummy := &http.Cookie{Name: "dummy", Value: "TS%3A1737311369%00%00csrf_token%3Adummy_csrf_token%00"}
	rSession := &http.Cookie{Name: "REVEL_SESSION", Value: "TS%3A1737311369%00%00csrf_token%3Atest_csrf_token%00"}
	noCSRF := &http.Cookie{Name: "REVEL_SESSION", Value: "TS%3A1737311369%00"}
	noTerminate := &http.Cookie{Name: "REVEL_SESSION", Value: "TS%3A1737311369%00%00csrf_token%3Atest_csrf_token"}
	invalidToken := &http.Cookie{Name: "REVEL_SESSION", Value: "TS%3A1737311369%00%00csrf_token%3A%test_csrf_token%00"}

	tests := []struct {
		name    string
		cookies []*http.Cookie
		want    string
	}{
		{
			name:    "valid",
			cookies: []*http.Cookie{dummy, rSession},
			want:    "test_csrf_token",
		},
		{
			name:    "no REVEL_SESSION",
			cookies: []*http.Cookie{dummy},
			want:    "",
		},
		{
			name:    "no csrf_token",
			cookies: []*http.Cookie{dummy, noCSRF},
			want:    "",
		},
		{
			name:    "no terminate",
			cookies: []*http.Cookie{dummy, noTerminate},
			want:    "",
		},
		{
			name:    "invalid token",
			cookies: []*http.Cookie{dummy, invalidToken},
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, csrf.FromCookies(context.Background(), tt.cookies))
		})
	}
}
