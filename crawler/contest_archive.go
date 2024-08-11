package crawler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/timezone"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/util"
	"github.com/pkg/errors"
)

var errNoPagenation = errors.New("no pagination list")

type ContestArchive struct {
	crawler *Crawler
}

func NewContestArchive(client *http.Client) *ContestArchive {
	return &ContestArchive{crawler: NewCrawler(url.ContestArchivePath).WithClient(client)}
}

func (c *ContestArchive) Do(ctx context.Context, req requests.ContestArchive) (*responses.ContestArchive, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	logger := logs.FromContext(ctx).With(
		slog.Group("request",
			slog.Any("relatedType", req.RatedType),
			slog.Any("page", req.Page),
		),
	)
	ctx = logs.ContextWith(ctx, logger)

	doc, err := c.crawler.DocumentGet(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to get document")
	}

	currentPage, err := c.parseCurrentPage(ctx, doc)
	if err != nil {
		if errors.Is(err, errNoPagenation) {
			logger.Info("no pagination list")
			return &responses.ContestArchive{}, nil
		}
		logger.Error(err.Error())
		return nil, errors.New("failed to parse current page")
	}

	totalPages, err := c.parseTotalPages(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse total pages")
	}

	contests, err := c.parseContents(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse contents")
	}

	return &responses.ContestArchive{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Contests:    contests,
	}, nil

}

func (c *ContestArchive) parseCurrentPage(ctx context.Context, doc *goquery.Document) (int, error) {
	cpText := doc.FindMatcher(goquery.Single("ul.pagination li.active a")).First().Text()
	logs.FromContext(ctx).With("text", cpText).Debug("find result")
	if len(cpText) == 0 {
		return 0, errNoPagenation
	}
	return strconv.Atoi(cpText)
}

func (c *ContestArchive) parseTotalPages(ctx context.Context, doc *goquery.Document) (int, error) {
	tpText := doc.Find("ul.pagination li").Last().Text()
	logs.FromContext(ctx).With("text", tpText).Debug("find result")
	return strconv.Atoi(tpText)
}

func (c *ContestArchive) parseContents(ctx context.Context, doc *goquery.Document) (responses.ContestArchive_ContestList, error) {
	logger := logs.FromContext(ctx)
	table := doc.Find("table.table.table-default").FilterFunction(func(i int, s *goquery.Selection) bool {
		if s.Find("table > thead tr th").FilterFunction(func(i int, s *goquery.Selection) bool {
			return s.Text() == "Contest Name"
		}).Length() == 0 {
			return false
		}
		if s.FindMatcher(goquery.Single("table > tbody tr td")).Length() == 0 {
			return false
		}
		return true
	})
	if table.Length() == 0 {
		logger.Debug("no contest list table is found")
		return nil, nil
	}
	trs := table.Find("table.table.table-default tbody tr")
	trsLen := trs.Length()
	var contests responses.ContestArchive_ContestList
	for i := range trsLen {
		logger := logger.With("index", i)
		c, err := c.parseContest(logs.ContextWith(ctx, logger), trs.Eq(i))
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to parse contest")
		}
		contests = append(contests, *c)
	}
	return contests, nil
}

func (c *ContestArchive) parseContest(ctx context.Context, tr *goquery.Selection) (*responses.ContestArchive_Contest, error) {
	logger := logs.FromContext(ctx)
	tds := tr.Find("td")

	// <td><a href="..."><time class="fixtime-full">2024-02-24(土) 21:00</time></a></td>
	// timeタグの中はgoqueryの解析では年月日時分秒+タイムゾーンの形式になってる?
	tdTime := tds.Eq(0).FindMatcher(goquery.Single("td a time.fixtime-full"))
	if tdTime.Length() == 0 {
		return nil, errors.New("no time is found")
	}
	startAt, err := time.ParseInLocation("2006-01-02 15:04:05-0700", tdTime.Text(), timezone.Tokyo)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse start time")
	}

	// <td><span/><span/><a href="/contests/abc342">HUAWEI Programming Contest 2024</a></td>
	td1Anchor := tds.Eq(1).Find("td a").FilterFunction(func(i int, s *goquery.Selection) bool {
		if href, ok := s.Attr("href"); ok && strings.HasPrefix(href, "/contests/") {
			return true
		}
		return false
	}).First()
	title := strings.TrimSpace(td1Anchor.Text())
	href, ok := td1Anchor.Attr("href")
	if !ok {
		logger.With("title", title).Error("no contest url is found")
		return nil, errors.New("failed to parse contest url")
	}
	id := ids.ContestID(strings.TrimPrefix(href, "/contests/"))
	if err := id.Validate(); err != nil {
		logger.With("id", id).Error(err.Error())
		return nil, errors.New("failed to validate contest id")
	}

	// <td>01:40</td>
	// 時間は3桁もありうる
	td2 := tds.Eq(2).Text()
	duration, err := util.ParseHoursMinutes(td2)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse duration")
	}

	// <td> ~ 1999</td>
	td3 := tds.Eq(3)
	targetRate := td3.Text()
	if targetRate == "-" {
		targetRate = "Unrated"
	}

	return &responses.ContestArchive_Contest{
		ID:         id,
		Title:      title,
		StartAt:    startAt,
		Duration:   duration,
		TargetRate: targetRate,
	}, nil
}
