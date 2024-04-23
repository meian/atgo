package crawler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/util"
	"github.com/pkg/errors"
)

type ContestTask struct {
	crawler *Crawler
}

func NewContestTask(client *http.Client) *ContestTask {
	return &ContestTask{crawler: NewCrawler(url.ContestTaskPath).WithClient(client)}
}

func (c *ContestTask) Do(ctx context.Context, req *requests.ContestTask) (*responses.ContestTask, error) {
	logger := logs.FromContext(ctx).With(
		slog.Group("request",
			slog.Any("contestID", req.ContestID),
		),
	)
	ctx = logs.ContextWith(ctx, logger)

	crawler := c.crawler.WithPathParam("contestID", req.ContestID)
	doc, err := crawler.DocumentGet(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to get document")
	}

	tasks, err := c.parseTasks(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse tasks")
	}

	return &responses.ContestTask{
		ContestID: req.ContestID,
		Tasks:     tasks,
	}, nil
}

func (c *ContestTask) parseTasks(ctx context.Context, doc *goquery.Document) ([]responses.ContestTask_Task, error) {
	logger := logs.FromContext(ctx)
	row := doc.Find("div.row").Find("div").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Find("h2").Text() == "Tasks"
	})
	if row.Length() == 0 {
		return nil, errors.New("failed to find tasks div")
	}
	trs := row.Find("table").Find("tbody").Find("tr")
	trLen := trs.Length()
	var tasks []responses.ContestTask_Task
	for i := range trLen {
		ctx := logs.ContextWith(ctx, logger.With("index", i))
		task, err := c.parseTask(ctx, trs.Eq(i))
		if err != nil {
			logs.FromContext(ctx).Error(err.Error())
			return nil, errors.Wrap(err, "failed to parse task")
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (c *ContestTask) parseTask(ctx context.Context, tr *goquery.Selection) (*responses.ContestTask_Task, error) {
	logger := logs.FromContext(ctx)
	// TODO: 個別メソッドにする
	tds := tr.Find("td")
	if tds.Length() < 4 {
		logger.Error("td must need more than 4")
		return nil, errors.New("failed to find td")
	}
	symbol := tds.Eq(0).Text()
	if len(symbol) == 0 {
		logger.Error("symbol must not be empty")
		return nil, errors.New("failed to find symbol")
	}
	title := tds.Eq(1).Text()
	if len(title) == 0 {
		logger.Error("title must not be empty")
		return nil, errors.New("failed to find title")
	}
	paths := strings.Split(tds.Eq(1).Find("a").AttrOr("href", ""), "/")
	if len(paths) == 0 {
		logger.Error("href must not be empty")
		return nil, errors.New("failed to find href")
	}
	id := paths[len(paths)-1]
	sec, err := util.ParseDuration(tds.Eq(2).Text())
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.Wrap(err, "failed to parse time limit")
	}
	mem, err := util.ParseMemory(tds.Eq(3).Text())
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.Wrap(err, "failed to parse memory limit")
	}
	return &responses.ContestTask_Task{
		ID:        id,
		Index:     symbol,
		Title:     title,
		TimeLimit: sec,
		Memory:    mem,
	}, nil
}
