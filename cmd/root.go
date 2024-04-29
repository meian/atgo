package cmd

import (
	"log/slog"
	gohttp "net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/meian/atgo/config"
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/flags"
	"github.com/meian/atgo/http"
	"github.com/meian/atgo/http/cookie"
	"github.com/meian/atgo/http/cookiestore"
	"github.com/meian/atgo/http/roundtrippers"
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/usecase"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootFlag struct {
	noLog    bool
	logDebug bool
	logInfo  bool
	logError bool
}

var rootCmd = &cobra.Command{
	Use:          "atgo",
	Short:        "A tool that provides an environment to participate in AtCoder for golang",
	Long:         `atgo is a tool that provides an environment to participate in AtCoder for golang.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.Inititalize()
		if err := initializeOutput(cmd); err != nil {
			return err
		}
		if err := initializeLogging(cmd); err != nil {
			return err
		}
		if err := initializeDatabase(cmd); err != nil {
			return err
		}
		if err := initializeHTTPClient(cmd); err != nil {
			return err
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		terminateHTTPClient(cmd)
	},
}

// initializeOutput は出力の初期化を行う
// コマンドの出力先をコンテキストに設定することで、コマンドと同じ出力先を他のレイヤーに共有できるようにしている
func initializeOutput(cmd *cobra.Command) error {
	ctx := io.OutWithContext(cmd.Context(), cmd.OutOrStdout())
	ctx = io.ErrWithContext(ctx, cmd.ErrOrStderr())
	cmd.SetContext(ctx)
	return nil
}

// initializeLogging はロガーの初期化を行う
// グローバルに設定されたフラグによって出力されるレベルを変更してロガーを作成する
func initializeLogging(cmd *cobra.Command) error {
	level, err := logs.ParseLevel(flags.DefaultLogLevel)
	if err != nil {
		return err
	}
	switch {
	case rootFlag.noLog:
		level = logs.LevelNone
	case rootFlag.logDebug:
		level = logs.LevelDebug
	case rootFlag.logInfo:
		level = logs.LevelInfo
	case rootFlag.logError:
		level = logs.LevelError
	}
	logger := logs.New(io.ErrFromContext(cmd.Context()), level)
	slog.SetDefault(logger)
	logger.Debug("setup logger")
	cmd.SetContext(logs.ContextWith(cmd.Context(), logger))
	return nil
}

// initializeDatabase はデータベース接続の初期化を行う
// DBがまだできていない場合は atgo init 相当の処理を事前に行うことで内部のデータを生成してから初期化を行う
// 一部コマンドでは呼び出された際にDBファイルを生成されると困るものもあるため、それらのコマンドではDBの初期化が行われないようになっている
// 具体的には database.Ignore() によって登録されたコマンドは初期化されず、対象コマンドの判別はコマンドのパスによって行っている
func initializeDatabase(cmd *cobra.Command) error {
	initFunc := func() error {
		_, err := usecase.Init{}.Run(cmd.Context(), usecase.InitParam{})
		return err
	}
	ctx, err := database.ContextWithNewDB(cmd.Context(), cmd.CommandPath(), initFunc)
	if err != nil {
		return err
	}
	cmd.SetContext(ctx)
	return nil
}

func initializeHTTPClient(cmd *cobra.Command) error {
	logger := logs.FromContext(cmd.Context())
	jopt := cookie.JarOption{
		IgnorePaths: []string{url.HomePath},
	}
	baseJar, _ := cookiejar.New(nil)
	jar := cookie.NewJar(baseJar, jopt)
	cfile, modtime, exists := workspace.CookieFile()
	logger = logger.With("cookie file", cfile).With("modtime", modtime)

	if exists {
		if time.Since(modtime) <= 24*time.Hour {
			url := url.URL("", nil, nil)
			if err := cookiestore.Load(url, cfile, jar); err != nil {
				logger.Error(err.Error())
				return errors.New("failed to load cookie")
			}
			logger.Info("load cookie")
		} else {
			if err := os.Remove(cfile); err != nil {
				logger.With("error", err.Error()).Warn("failed to remove old cookie file")
			}
			logger.Info("skip loading because cookie file is too old")
		}
	}
	rt := roundtrippers.NewRateLimitRoundTripper(gohttp.DefaultTransport, 1*time.Second)
	client := http.NewClient(rt, jar)
	cmd.SetContext(http.ContextWithClient(cmd.Context(), client))
	return nil
}

func terminateHTTPClient(cmd *cobra.Command) {
	logger := logs.FromContext(cmd.Context())
	client := http.ClientFromContext(cmd.Context())
	if client == nil || client.Jar == nil {
		logger.Debug("skip store cookie because cookie jar is not found")
		return
	}
	cfile, _, _ := workspace.CookieFile()
	logger = logger.With("cookie file", cfile)
	if err := cookiestore.Save(url.URL("", nil, nil), cfile, client.Jar); err != nil {
		logger.With("error", err.Error()).Warn("failed to store cookies")
		return
	}
	logger.Info("store cookies")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootFlag.noLog, "nolog", "n", false, "no output log")
	rootCmd.PersistentFlags().BoolVarP(&rootFlag.logDebug, "verbose", "v", false, "output log specified by error/warn/info/debug")
	rootCmd.PersistentFlags().BoolVarP(&rootFlag.logInfo, "info", "i", false, "output log specified by error/warn/info")
	rootCmd.PersistentFlags().BoolVarP(&rootFlag.logError, "error", "e", false, "output log only error")

	rootCmd.SetErrPrefix(color.New(color.FgHiRed).SprintFunc()("ERR"))
}
