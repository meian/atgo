package cmd

import (
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/tmpl"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var initFlag struct {
	force bool
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize atgo managed database",
	Long: `Initialize atgo managed database.
This command contains
- Create database file
- Create tables
- Insert initial data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		p := usecase.InitParam{
			Force: initFlag.force,
		}
		res, err := usecase.Init{}.Run(ctx, p)
		if err != nil {
			return err
		}

		tname := "init"
		t := tmpl.CmdTemplate(tname)
		w := io.OutFromContext(ctx)
		if err := t.Execute(w, res); err != nil {
			logs.FromContext(ctx).With("template name", tname).Error(err.Error())
			return errors.New("failed to execute template")
		}

		return nil
	},
}

func init() {
	database.Ignore(initCmd.CommandPath)
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initFlag.force, "force", "f", false, "if db is already initialized, clean up and reinitialize")
}
