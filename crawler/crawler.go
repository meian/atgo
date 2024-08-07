package crawler

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Crawler struct {
	Path       string
	pathParams map[string]string
	client     *http.Client
}

func NewCrawler(path string) *Crawler {
	return &Crawler{
		Path:       path,
		pathParams: make(map[string]string),
	}
}

func (c *Crawler) WithClient(client *http.Client) *Crawler {
	c = c.clone()
	c.client = client
	return c
}

func (c *Crawler) clone() *Crawler {
	ppm := make(map[string]string, len(c.pathParams))
	for k, v := range c.pathParams {
		ppm[k] = v
	}
	return &Crawler{
		Path:       c.Path,
		pathParams: ppm,
		client:     c.client,
	}
}

func (c *Crawler) WithPathParam(key, value string) *Crawler {
	c = c.clone()
	c.pathParams[key] = value
	return c
}

func (c Crawler) Get(ctx context.Context, queries url.Valuer) (*http.Response, error) {
	return c.crawl(ctx, "GET", "", queries, nil)
}

func (c Crawler) DocumentGet(ctx context.Context, queries url.Valuer) (*goquery.Document, error) {
	resp, err := c.Get(ctx, queries)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.documentFromReader(ctx, resp.Body)
}

func (c Crawler) Post(ctx context.Context, queries, bodies url.Valuer) (*http.Response, error) {
	return c.crawl(ctx, "POST", "application/x-www-form-urlencoded", queries, bodies)
}

func (c Crawler) DocumentPost(ctx context.Context, queries, bodies url.Valuer) (*goquery.Document, error) {
	resp, err := c.Post(ctx, queries, bodies)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.documentFromReader(ctx, resp.Body)
}

func (c Crawler) crawl(ctx context.Context, method, contentType string, queries, bodies url.Valuer) (*http.Response, error) {
	url := url.URL(c.Path, c.pathParams, queries)
	logger := logs.FromContext(ctx).
		With("url", url.String()).
		With("method", method)
	ctx = logs.ContextWith(ctx, logger)

	ctx, body := c.bodyReader(ctx, bodies)
	return c.doHTTPRequest(ctx, method, contentType, url.String(), body)
}

func (c Crawler) bodyReader(ctx context.Context, bodies url.Valuer) (context.Context, io.Reader) {
	if bodies == nil {
		logs.FromContext(ctx).Debug("no bodies")
		return ctx, nil
	}

	logger := logs.FromContext(ctx)
	logger.Debug("creating body reader")
	bodyText := bodies.URLValues().Encode()
	logger = logger.With(
		slog.Group("body",
			slog.Any("count", strings.Count(bodyText, "&")+1),
			slog.Any("length", len(bodyText)),
		),
	)
	body := strings.NewReader(bodyText)
	logger.Debug("created body reader")

	return logs.ContextWith(ctx, logger), body
}

func (c Crawler) doHTTPRequest(ctx context.Context, method, contentType, url string, body io.Reader) (*http.Response, error) {
	logger := logs.FromContext(ctx)
	logger.Debug("creating request")
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to create request")
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	logger.Debug("created request")

	logger.Debug("sending request")
	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to send request")
	}
	logger.Debug("sent request")

	if resp.StatusCode >= http.StatusBadRequest {
		logger.With("statusCode", resp.StatusCode).Error("status code is not 200")
		return nil, errors.New("status code is not 200")
	}

	return resp, nil
}

func (c Crawler) documentFromReader(ctx context.Context, respBody io.Reader) (*goquery.Document, error) {
	logger := logs.FromContext(ctx)
	logger.Debug("parsing document from response")

	respBody, err := io.WithReadAction(respBody, func(r io.Reader) error {
		return c.validHTML(ctx, r)
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("response is not a valid HTML")
	}

	doc, err := goquery.NewDocumentFromReader(respBody)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse document from response")
	}
	logger.Debug("parsed document from response")

	return doc, nil
}

func (Crawler) validHTML(ctx context.Context, reader io.Reader) error {
	logger := logs.FromContext(ctx)

	tags := [][]byte{[]byte("<html"), []byte("<head"), []byte("<body")}
	buf := make([]byte, 4096)
	var content []byte

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				return errors.New("failed to find HTML tags")
			}
			logger.Error(err.Error())
			return errors.New("failed to read response body")
		}
		content = append(content, buf[:n]...)

		for {
			pos := bytes.Index(content, tags[0])
			if pos < 0 {
				break
			}
			content = content[pos:]
			tags = tags[1:]
			if len(tags) == 0 {
				return nil
			}
		}
	}
}

func (c Crawler) LoggedIn(ctx context.Context, doc *goquery.Document) bool {
	selector := fmt.Sprintf("a[href='%s']", url.SettingsPath)
	return doc.Find(selector).Length() > 0
}
