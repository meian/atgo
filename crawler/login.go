package crawler

import (
	"context"
	"net/http"
	gourl "net/url"

	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/csrf"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Login struct {
	crawler *Crawler
}

func NewLogin(client *http.Client) *Login {
	crawler := NewCrawler(url.LoginPath).
		WithClient(client).
		WithLoadCookie(false)
	return &Login{crawler: crawler}
}

type loginQuery struct {
	Continue string
}

func (q *loginQuery) URLValues() gourl.Values {
	return gourl.Values{
		"continue": {q.Continue},
	}
}

func (c *Login) Do(ctx context.Context, req *requests.Login) (*responses.Login, error) {
	logger := logs.FromContext(ctx).With("continue", req.Continue)
	resp, err := c.crawler.Post(ctx, &loginQuery{
		Continue: req.Continue,
	}, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to post document")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.With("status_code", resp.StatusCode).Error("unexpected status code")
		return nil, errors.New("unexpected status code for login")
	}
	return &responses.Login{
		LoggedIn:  resp.Request.URL.String() == req.Continue,
		CSRFToken: csrf.FromCookies(ctx, resp.Cookies()),
	}, nil
}
