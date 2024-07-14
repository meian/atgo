package crawler

import (
	"context"
	gohttp "net/http"
	gourl "net/url"

	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/csrf"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Login struct {
	crawler *Crawler
}

func NewLogin(client *gohttp.Client) *Login {
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
	ctx = http.ContextWithSkipWait(ctx)
	resp, err := c.crawler.Post(ctx, &loginQuery{
		Continue: req.Continue,
	}, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to post document")
	}
	defer resp.Body.Close()
	if resp.StatusCode != gohttp.StatusOK {
		logger.With("statusCode", resp.StatusCode).Error("unexpected status code")
		return nil, errors.New("unexpected status code for login")
	}
	doc, err := c.crawler.documentFromReader(ctx, resp.Body)
	if err != nil {
		return nil, err
	}
	return &responses.Login{
		LoggedIn:  c.crawler.LoggedIn(ctx, doc),
		CSRFToken: csrf.FromCookies(ctx, resp.Cookies()),
	}, nil
}
