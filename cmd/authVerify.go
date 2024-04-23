package cmd

import (
	"github.com/meian/atgo/auth"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var authVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify stored credential",
	Long:  `Verify stored credential in your local environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file := workspace.CredentialFile()
		ctx := cmd.Context()
		logger := logs.FromContext(ctx)
		if !io.FileExists(file) {
			return errors.New("No credential is stored.")
		}
		user, password, err := auth.Read(ctx, file)
		if err != nil {
			logger.Error(err.Error())
			return errors.New("failed to read credential")
		}
		if err := tryLogin(ctx, user, password); err != nil {
			return err
		}
		cmd.Printf("Credential is verified: %s\n", file)
		cmd.Printf("Username: %s\n", user)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authVerifyCmd)
}
