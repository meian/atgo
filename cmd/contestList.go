package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/help"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/text"
	"github.com/meian/atgo/tmpl"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var contestListFlag struct {
	size int
	page int
}

var contestListCmd = &cobra.Command{
	Use:   "list [abc/arc/agc/ahc or all]",
	Short: "Display a list of locally stored contests",
	Long: `Display a list of locally stored contests.
The list is retrieved from the following commands that have been saved.

* contest list retrieved by ` + "`{{ .load }}`." + `
* a contest retrieved by ` + "`{{ .contest }}`.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if len(args) == 0 {
			cmd.Usage()
			return errors.New("rated type is required")
		}
		ratedType := args[0]
		logger := logs.FromContext(ctx).With("ratedType", ratedType)
		ctx = logs.ContextWith(ctx, logger)
		if !slices.Contains(constant.RatedTypeStrings(), ratedType) {
			cmd.Usage()
			return errors.Errorf("invalid rated type: %s", ratedType)
		}

		p := usecase.ContestListParam{
			RatedType: ratedType,
			Page:      contestListFlag.page,
			Size:      contestListFlag.size,
		}
		res, err := usecase.ContestList{}.Run(ctx, p)
		if err != nil {
			return err
		}

		tname := "contest_list"
		t := tmpl.CmdTemplate(tname)
		w := io.OutFromContext(ctx)
		model := toContestListOutputModel(res, ratedType)
		if err := t.Execute(w, model); err != nil {
			logger.With("template name", tname).Error(err.Error())
			return errors.New("failed to execute template")
		}

		return nil
	},
}

func init() {
	help.Replace(contestListCmd, []help.Replacer{
		{Key: "load", Command: contestLoadCmd},
		{Key: "contest", Command: contestCmd},
	})
	contestCmd.AddCommand(contestListCmd)

	contestListCmd.Flags().IntVarP(&contestListFlag.size, "size", "s", 20, "size of show contest list")
	contestListCmd.Flags().IntVarP(&contestListFlag.page, "page", "p", 1, "page number in contest list")
}

func toContestListOutputModel(res *usecase.ContestListResult, ratedType string) *ContestListOutputModel {
	idLen, titleLen := 0, 0
	var contests []ContestListOutputModel_Contest
	for _, c := range res.Contests {
		idLen = max(idLen, text.StringWidth(c.ID))
		titleLen = max(titleLen, text.StringWidth(c.Title))
		contests = append(contests, ContestListOutputModel_Contest{
			ID:      c.ID,
			Title:   c.Title,
			StartAt: c.StartAt,
			EndAt:   c.StartAt.Add(c.Duration),
		})
	}
	for i := range contests {
		contests[i].IDPadding = text.PadRight(contests[i].ID, idLen)
		contests[i].TitlePadding = text.PadRight(contests[i].Title, titleLen)
	}
	return &ContestListOutputModel{
		LoadCommand: fmt.Sprintf("%s %s", contestLoadCmd.CommandPath(), ratedType),
		Contests:    contests,
	}
}

type ContestListOutputModel struct {
	LoadCommand string
	Contests    []ContestListOutputModel_Contest
}

type ContestListOutputModel_Contest struct {
	ID           string
	IDPadding    string
	Title        string
	TitlePadding string
	StartAt      time.Time
	EndAt        time.Time
}
