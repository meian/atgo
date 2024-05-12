package cmd

import (
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/flags"
	"github.com/meian/atgo/http/cookie"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/tmpl"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var versionFlag struct {
	Long bool
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show app version",
	Long:  `Show app version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !versionFlag.Long {
			v := flags.Version
			if v == "" {
				v = "no_version"
			}
			cmd.Println(v)
			return nil
		}
		model := VersionOutputModel{
			Version:     flags.Version,
			CommitSHA:   flags.CommitSHA,
			Description: rootCmd.Short,
			FlagsPkg:    flags.Package(),
		}

		ctx := cmd.Context()
		logger := logs.FromContext(ctx)

		tname := "version"
		t := tmpl.CmdTemplate(tname)
		w := io.OutFromContext(ctx)
		if err := t.Execute(w, model); err != nil {
			logger.With("template name", tname).Error(err.Error())
			return errors.New("failed to execute template")
		}
		return nil
	},
}

func init() {
	cookie.IgnoreLoad(versionCmd.CommandPath)
	cookie.IgnoreSave(versionCmd.CommandPath)
	database.Ignore(versionCmd.CommandPath)
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&versionFlag.Long, "long", false, "Show version and descriptions")
}

type VersionOutputModel struct {
	Version     string
	CommitSHA   string
	Description string
	FlagsPkg    string
}
