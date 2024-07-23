package csrf

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/meian/atgo/logs"
)

func FromCookies(ctx context.Context, cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "REVEL_SESSION" {
			return fromCookie(ctx, cookie)
		}
	}
	return ""
}

func fromCookie(ctx context.Context, cookie *http.Cookie) string {
	logger := logs.FromContext(ctx)
	const header = "csrf_token%3A"
	tokenPos := strings.Index(cookie.Value, header)
	if tokenPos < 0 {
		logger.Warn("csrf token not found in cookie")
		return ""
	}
	v := cookie.Value[tokenPos+len(header):]
	endPos := strings.Index(v, "%00")
	if endPos < 0 {
		logger.Error("csrf token is not terminate with %00")
		return ""
	}
	v = v[:endPos]
	v, err := url.QueryUnescape(v)
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	return v
}
