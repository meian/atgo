package cmd

import (
	"github.com/meian/atgo/help"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var submitFlag struct {
	skipTest bool
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your code to AtCoder",
	Long: `Submit your code saved in local main.go to AtCoder.
Code is sent to the task set the last locally with ` + "`{{ .localInit }}`" + `
Before submitting code, it will attempt to build main.go and run unit testing under main package,
and it is skipped if either are failed.
For example, if there are multiple patterns of correct output, and you want to send even if test is failed,
you can skip test and send by adding --skip-test.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := usecase.SubmitParam{
			SkipTest: submitFlag.skipTest,
		}
		res, err := usecase.Submit{}.Run(cmd.Context(), p)
		if err != nil {
			return err
		}
		switch res.ErrorStage {
		case usecase.SubmitErrorStageBuild:
			return errors.Errorf("failed to build\n%s", res.ErrorMessage)
		case usecase.SubmitErrorStageTest:
			return errors.Errorf("failed to test\n%s\nfor skip test, you can use `%s --skip-test`", res.ErrorMessage, cmd.CommandPath())
		}
		if !res.LoggedIn {
			return errors.New("failed to login")
		}
		if !res.Submitted {
			return errors.New("failed to submit")
		}
		cmd.Println("complete to submit")
		return nil
	},
}

func init() {
	help.Replace(submitCmd, []help.Replacer{
		{Key: "localInit", Command: taskLocalInitCmd},
	})
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVarP(&submitFlag.skipTest, "skip-test", "s", false, "skip running unit test before submitting")
}
