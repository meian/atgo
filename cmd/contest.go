package cmd

import (
	"fmt"
	"time"

	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/text"
	"github.com/meian/atgo/tmpl"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var contestCmd = &cobra.Command{
	Use:   "contest [contest ID]",
	Short: "Display contest info",
	Long:  `Display contest and associated tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var contestID string
		if len(args) > 0 {
			contestID = args[0]
		}
		logger := logs.FromContext(ctx).With("contestID", contestID)
		ctx = logs.ContextWith(ctx, logger)

		p := usecase.ContestParam{
			ContestID: contestID,
		}
		res, err := usecase.Contest{}.Run(ctx, p)
		if err != nil {
			return err
		}

		tname := "contest"
		t := tmpl.CmdTemplate(tname)
		w := io.OutFromContext(ctx)
		model := toContestOutputModel(res)
		if err := t.Execute(w, model); err != nil {
			logger.With("template name", tname).Error(err.Error())
			return errors.New("failed to execute template")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(contestCmd)
}

func toContestOutputModel(res *usecase.ContestResult) ContestOutputModel {
	contest := res.Contest
	var ratedType string
	if contest.RatedType.Valid {
		ratedType = contest.RatedType.String
	} else {
		ratedType = "no rate info"
	}
	tasks := make([]ContestOutputModel_Task, len(contest.ContestTasks))
	indexLen, idLen, titleLen, scoreLen := 0, 0, 0, 0
	for i, ct := range contest.ContestTasks {
		indexLen = max(indexLen, text.StringWidth(ct.Index))
		idLen = max(idLen, text.StringWidth(ct.Task.ID))
		titleLen = max(titleLen, text.StringWidth(ct.Task.Title))
		var score string
		if ct.Task.Score.Valid {
			score = fmt.Sprintf("%d pt", ct.Task.Score.Int64)
		} else {
			score = "no score info"
		}
		scoreLen = max(scoreLen, text.StringWidth(score))

		tasks[i] = ContestOutputModel_Task{
			Index:     ct.Index,
			ID:        ct.TaskID,
			Title:     ct.Task.Title,
			Score:     score,
			TimeLimit: ct.Task.TimeLimit,
			Memory:    ct.Task.Memory,
		}
	}
	for i := range tasks {
		tasks[i].IndexPadding = text.PadRight(tasks[i].Index, indexLen)
		tasks[i].IDPadding = text.PadRight(tasks[i].ID, idLen)
		tasks[i].TitlePadding = text.PadRight(tasks[i].Title, titleLen)
		tasks[i].ScorePadding = text.PadLeft(tasks[i].Score, scoreLen)
	}
	return ContestOutputModel{
		ID:         contest.ID,
		RatedType:  ratedType,
		Title:      contest.Title,
		StartAt:    contest.StartAt,
		EndAt:      contest.StartAt.Add(contest.Duration),
		TargetRate: contest.TargetRate,
		Tasks:      tasks,
	}

}

type ContestOutputModel struct {
	ID         string
	RatedType  string
	Title      string
	StartAt    time.Time
	EndAt      time.Time
	TargetRate string
	Tasks      []ContestOutputModel_Task
}

type ContestOutputModel_Task struct {
	Index        string
	IndexPadding string
	ID           string
	IDPadding    string
	Title        string
	TitlePadding string
	Score        string
	ScorePadding string
	TimeLimit    time.Duration
	Memory       int
}
