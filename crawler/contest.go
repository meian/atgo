package crawler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/timezone"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Contest struct {
	crawler *Crawler
}

func NewContest(client *http.Client) *Contest {
	return &Contest{crawler: NewCrawler(url.ContestPath).WithClient(client)}
}

func (c *Contest) Do(ctx context.Context, req requests.Contest) (*responses.Contest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	logger := logs.FromContext(ctx).With(
		slog.Group("request",
			slog.Any("contestID", req.ContestID),
		),
	)
	ctx = logs.ContextWith(ctx, logger)

	crawler := c.crawler.WithPathParam("contestID", string(req.ContestID))
	doc, err := crawler.DocumentGet(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to get document")
	}

	title, err := c.parseTitle(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse title")
	}

	startAt, duration, err := c.parseTimes(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse times")
	}

	targetRate, err := c.parseTargetRate(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse target rate")
	}

	return &responses.Contest{
		ID:         req.ContestID,
		Title:      title,
		StartAt:    startAt,
		Duration:   duration,
		TargetRate: targetRate,
	}, nil
}

func (c *Contest) parseTitle(ctx context.Context, doc *goquery.Document) (string, error) {
	// <nav><ul class="nav"><li><a class="contest-title" href="/contests/abc342">AtCoder Beginner Contest 342</a></li></ul></nav>
	title := strings.TrimSpace(doc.FindMatcher(goquery.Single("nav ul.nav li a.contest-title")).Text())
	logs.FromContext(ctx).With("text", title).Debug("find result")
	if title == "" {
		return "", errors.New("no title is found")
	}
	return title, nil
}

func (c *Contest) parseTimes(ctx context.Context, doc *goquery.Document) (time.Time, time.Duration, error) {
	// <small class="contest-duration"><time class="fixtime-full">2024-02-24(土) 21:00</time></small>
	// timeタグが二つあって開始〜終了の時間が書いてある
	tt := doc.Find("small.contest-duration time.fixtime-full")
	if tt.Length() != 2 {
		return time.Time{}, 0, errors.New("no time is found")
	}
	logs.FromContext(ctx).
		With("startAt", tt.Eq(0).Text()).
		With("endAt", tt.Eq(1).Text()).
		Debug("find result")
	startAt, err := time.ParseInLocation("2006-01-02 15:04:05-0700", tt.Eq(0).Text(), timezone.Tokyo)
	if err != nil {
		return time.Time{}, 0, errors.New("failed to parse start time")
	}
	endAt, err := time.ParseInLocation("2006-01-02 15:04:05-0700", tt.Eq(1).Text(), timezone.Tokyo)
	if err != nil {
		return time.Time{}, 0, errors.New("failed to parse end time")
	}
	return startAt, endAt.Sub(startAt), nil
}

func (c *Contest) parseTargetRate(ctx context.Context, doc *goquery.Document) (string, error) {
	// <div id="main-container"><div class="row"><div><p class="small"><span>Rated Range:  ~ 1999</span></p></div></div></div>
	tr := doc.Find("div#main-container div.row p.small span").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.HasPrefix(s.Text(), "Rated Range: ")
	}).Text()
	logs.FromContext(ctx).With("text", tr).Debug("find result")
	if len(tr) == 0 {
		return "", errors.New("no target rate is found")
	}
	return strings.TrimPrefix(tr, "Rated Range: "), nil
}
