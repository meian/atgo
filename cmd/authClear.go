package cmd

import (
	"os"

	"github.com/meian/atgo/io"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var authClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear stored credential",
	Long:  `Clear stored credential in your local environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file := workspace.CredentialFile()
		if !io.FileExists(file) {
			cmd.Println("No credential is stored.")
			return nil
		}
		if err := os.Remove(file); err != nil {
			return errors.Wrap(err, "failed to remove credential file")
		}
		cmd.Println("Credential is removed.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authClearCmd)
}
