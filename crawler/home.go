package crawler

import (
	"context"
	gohttp "net/http"

	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/csrf"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Home struct {
	crawler *Crawler
}

func NewHome(client *gohttp.Client) *Home {
	crawler := NewCrawler(url.HomePath).WithClient(client)
	return &Home{crawler: crawler}
}

func (c *Home) Do(ctx context.Context, req *requests.Home) (*responses.Home, error) {
	ctx = http.ContextWithSkipWait(ctx)
	doc, err := c.crawler.DocumentGet(ctx, req)
	if err != nil {
		logs.FromContext(ctx).Error(err.Error())
		return nil, errors.New("failed to get document")
	}
	return &responses.Home{
		CSRFToken: csrf.FromDocument(doc),
	}, nil
}
