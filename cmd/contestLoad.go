package cmd

import (
	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/help"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/tmpl"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var contestLoadCmd = &cobra.Command{
	Use:   "load [abc/arc/agc/ahc]",
	Short: "Load archived contests and store it locally",
	Long: `Retrieve a list of previously held contests on AtCoder and store it locally.
The saved contests are available by ` + "`{{ .list }}.`",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if len(args) == 0 {
			cmd.Usage()
			return errors.New("rated type is required")
		}
		rt := args[0]
		logger := logs.FromContext(ctx).With("ratedType", rt)
		ctx = logs.ContextWith(ctx, logger)
		ratedType, err := constant.RatedTypeString(rt)
		if err != nil || ratedType == constant.RatedTypeAll {
			cmd.Usage()
			return errors.Errorf("invalid rated type: %s", rt)
		}

		tname := "contest_load"
		t := tmpl.CmdTemplate(tname)
		w := io.OutFromContext(ctx)

		for cp, tp := 1, 10000; ; cp++ {
			p := usecase.ContestLoadParam{
				RatedType: ratedType,
				Page:      cp,
			}
			res, err := usecase.ContestLoad{}.Run(ctx, p)
			if err != nil {
				return err
			}
			if err := t.Execute(w, res); err != nil {
				logger.With("template name", tname).Error(err.Error())
				return errors.New("failed to execute template")
			}
			tp = res.TotalPages
			if cp >= tp {
				break
			}
		}

		return nil
	},
}

func init() {
	help.Replace(contestLoadCmd, []help.Replacer{
		{Key: "list", Command: contestListCmd},
	})
	contestCmd.AddCommand(contestLoadCmd)

}
