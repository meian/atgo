package crawler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/crawler/responses"
	"github.com/meian/atgo/csrf"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/url"
	"github.com/pkg/errors"
)

type Task struct {
	crawler *Crawler
}

func NewTask(client *http.Client) *Task {
	return &Task{crawler: NewCrawler(url.TaskPath).WithClient(client)}
}

func (c *Task) Do(ctx context.Context, req requests.Task) (*responses.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	logger := logs.FromContext(ctx).With(
		slog.Group("request",
			slog.Any("contestID", req.ContestID),
			slog.Any("taskID", req.TaskID),
		),
	)
	ctx = logs.ContextWith(ctx, logger)

	crawler := c.crawler.
		WithPathParam("contestID", req.ContestID).
		WithPathParam("id", req.TaskID)
	doc, err := crawler.DocumentGet(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to get document")
	}

	score, err := c.parseScore(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse score")
	}

	samples, err := c.parseSamples(ctx, doc)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse samples")
	}

	return &responses.Task{
		ID:        ids.TaskID(req.TaskID),
		Score:     score,
		Samples:   samples,
		CSRFToken: csrf.FromDocument(doc),
		LoggedIn:  c.crawler.LoggedIn(ctx, doc),
	}, nil
}

func (c *Task) parseScore(ctx context.Context, doc *goquery.Document) (*int, error) {

	stmt, err := c.taskStatement(ctx, doc)
	if err != nil {
		return nil, err
	}

	logger := logs.FromContext(ctx)
	// 実際のHTMLでは <span.lang-ja><p><span.katex-mathml><math><semantics><annotation> の深い箇所にあるが、
	// goquery では <span.lang-ja><p><var> の階層で取得できる
	scorePart := stmt.Find("span.lang-ja > p").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(s.Text(), "配点")
	}).First().Find("p > var")
	if scorePart.Length() == 0 {
		logger.Debug("not found score part")
		return nil, nil
	}
	text := scorePart.Text()

	score, err := strconv.Atoi(text)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("failed to parse task score")
	}
	return &score, nil
}

func (c *Task) parseSamples(ctx context.Context, doc *goquery.Document) ([]responses.Task_Sample, error) {

	stmt, err := c.taskStatement(ctx, doc)
	if err != nil {
		return nil, err
	}

	// Sample Input の数 = sampleの数
	count := doc.Find("h3").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.HasPrefix(s.Text(), "入力例")
	}).Length()
	logger := logs.FromContext(ctx).With("sample count", count)

	var samples []responses.Task_Sample
	for i := 0; i < count; i++ {
		logger := logger.With("index", i)
		// 中に 入力例 N がある<div class="part">の中の<pre>を取得
		inTitle := fmt.Sprintf("入力例%d", i+1)
		input := stmt.Find("div.part").FilterFunction(func(i int, s *goquery.Selection) bool {
			return strings.ReplaceAll(s.Find("h3").Text(), " ", "") == inTitle
		}).Find("pre").First()
		if input.Length() == 0 {
			logger.Error("not found input")
			return nil, errors.New("failed to find input")
		}
		// 中に 出力例 N がある<div class="part">の中の<pre>を取得
		outTitle := fmt.Sprintf("出力例%d", i+1)
		output := stmt.Find("div.part").FilterFunction(func(i int, s *goquery.Selection) bool {
			return strings.ReplaceAll(s.Find("h3").Text(), " ", "") == outTitle
		}).Find("pre").First()
		if output.Length() == 0 {
			logger.Error("not found output")
			return nil, errors.New("failed to find output")
		}
		samples = append(samples, responses.Task_Sample{
			Input:  input.Text(),
			Output: output.Text(),
		})
	}
	return samples, nil
}

func (c *Task) taskStatement(_ context.Context, doc *goquery.Document) (*goquery.Selection, error) {
	// span.stmt がない場合もあるので、その場合は div.task-statement を対象にする
	stmt := doc.Find("span.lang-ja")
	if stmt.Length() == 0 {
		stmt = doc.Find("div#task-statement").First()
		if stmt.Length() == 0 {
			return nil, errors.New("failed to find task statement")
		}
	}
	return stmt, nil
}
