package cmd

import (
	"os"

	"github.com/meian/atgo/auth"
	"github.com/spf13/cobra"
)

var authSampleFlag struct {
	user     string
	password string
}

// authSampleCmd represents the authSample command
var authSampleCmd = &cobra.Command{
	Use:    "auth-sample",
	Short:  "For development",
	Long:   `Write sample credential file`,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		tmp, err := os.CreateTemp("", "")
		if err != nil {
			return err
		}
		if err := tmp.Close(); err != nil {
			return err
		}
		file := tmp.Name()
		defer os.Remove(file)
		if err := auth.Write(cmd.Context(), file, authSampleFlag.user, authSampleFlag.password); err != nil {
			return err
		}
		c, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(c)
		return err
	},
}

func init() {
	rootCmd.AddCommand(authSampleCmd)
	authSampleCmd.Flags().StringVar(&authSampleFlag.user, "user", "sample-user", "User name")
	authSampleCmd.Flags().StringVar(&authSampleFlag.password, "password", "sample-password", "Password")
}
