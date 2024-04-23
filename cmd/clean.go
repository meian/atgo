package cmd

import (
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up atgo database file",
	Long:  `Clean up atgo database file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := usecase.Clean{}.Run(cmd.Context())
		if err != nil {
			return errors.New("failed to delete database file")
		}

		if res.AlreadyCleaned {
			cmd.Println("no database is created")
		} else if res.Cleaned {
			cmd.Println("cleaned up database")
		}

		return nil
	},
}

func init() {
	database.Ignore(cleanCmd.CommandPath)
	rootCmd.AddCommand(cleanCmd)
}
