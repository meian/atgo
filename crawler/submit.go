package crawler

import (
	"context"
	"net/http"

	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Submit struct {
	crawler *Crawler
}

func NewSubmit(client *http.Client) *Submit {
	return &Submit{crawler: NewCrawler(url.SubmitPath).WithClient(client)}
}

func (c *Submit) Do(ctx context.Context, req *requests.Submit) (*responses.Submit, error) {
	logger := logs.FromContext(ctx)
	crawler := c.crawler.WithPathParam("contestID", req.ContestID)
	resp, err := crawler.Post(ctx, nil, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to post document")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.With("statusCode", resp.StatusCode).Error("unexpected status code")
		return nil, errors.New("unexpected status code for submit")
	}
	// TODO: URLが違う場合にログを出す
	return &responses.Submit{
		Submitted: resp.Request.URL.String() == url.MySubmissionURL(ids.ContestID(req.ContestID)),
	}, nil
}
