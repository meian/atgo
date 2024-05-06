package cmd

import (
	"context"

	"github.com/manifoldco/promptui"
	"github.com/meian/atgo/auth"
	"github.com/meian/atgo/crawler"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/http/cookie"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Store your AtCoder credential",
	Long:  `Store your AtCoder credential in your local environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logs.FromContext(ctx)
		userp := promptui.Prompt{
			Label: "Username",
			Validate: func(s string) error {
				if s == "" {
					return errors.New("username is required")
				}
				return nil
			},
		}
		user, err := userp.Run()
		if err != nil {
			logger.Error(err.Error())
			return errors.New("failed to get username")
		}
		passwordp := promptui.Prompt{
			Label: "Password",
			Mask:  ' ',
			Validate: func(s string) error {
				if s == "" {
					return errors.New("password is required")
				}
				return nil
			},
		}
		password, err := passwordp.Run()
		if err != nil {
			logger.Error(err.Error())
			return errors.New("failed to get password")
		}
		if err := tryLogin(ctx, user, password); err != nil {
			return err
		}
		file := workspace.CredentialFile()
		if err := auth.Store(ctx, file, user, password); err != nil {
			return err
		}
		cmd.Printf("Credential is stored: %s\n", file)
		return nil
	},
}

func tryLogin(ctx context.Context, username, password string) error {
	logger := logs.FromContext(ctx)
	client := http.ClientFromContext(ctx)
	hres, err := crawler.NewHome(client).Do(ctx, &requests.Home{})
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to get home response")
	}
	req := &requests.Login{
		Username:  username,
		Password:  password,
		CSRFToken: hres.CSRFToken,
		Continue:  url.HomeURL(),
	}
	logger.With("len(username)", len(username)).
		With("len(password)", len(password)).
		With("len(csrfToken)", len(hres.CSRFToken)).
		Info("login param")
	res, err := crawler.NewLogin(client).Do(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to login response")
	}
	if !res.LoggedIn {
		return errors.New("failed to login")
	}
	return nil
}

func init() {
	cookie.IgnoreLoad(authCmd.CommandPath)
	rootCmd.AddCommand(authCmd)
}
