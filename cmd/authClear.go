package cmd

import (
	"fmt"
	"os"

	"github.com/meian/atgo/http/cookie"
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
		credError := removeCredential(cmd)
		cookieError := removeCookie(cmd)
		if credError != nil {
			if cookieError != nil {
				return fmt.Errorf("%w / %w", credError, cookieError)
			}
			return credError
		}
		return cookieError
	},
}

// TODO: usecase化
func removeCredential(cmd *cobra.Command) error {
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
}

func removeCookie(cmd *cobra.Command) error {
	file, _, exists := workspace.CookieFile()
	if !exists {
		cmd.Println("No cookie is stored.")
		return nil
	}
	if err := os.Remove(file); err != nil {
		return errors.Wrap(err, "failed to remove cookie file")
	}
	cmd.Println("Cookie is removed.")
	return nil

}

func init() {
	cookie.IgnoreSave(authClearCmd.CommandPath)
	authCmd.AddCommand(authClearCmd)
}
