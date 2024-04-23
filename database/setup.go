package database

import (
	"context"

	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/workspace"
	"github.com/pkg/errors"
)

var ignoreTargets []func() string

func ContextWithNewDB(ctx context.Context, commandPath string, createDBFunc func() error) (context.Context, error) {
	logger := logs.FromContext(ctx)
	for _, cpFunc := range ignoreTargets {
		if cpFunc() == commandPath {
			logger.With("commandPath", commandPath).Info("ignore setup database")
			return ctx, nil
		}
	}
	logger.Info("setup database connection")
	_, exists := workspace.DBFile()
	if !exists {
		logger.Debug("not found database file, so create new database")
		if err := createDBFunc(); err != nil {
			logger.Error(err.Error())
			return nil, errors.New("failed to create database")
		}
	}
	return ContextWith(ctx, New(ctx)), nil
}

// Ignore は指定されたコマンドパスの場合にはデータベースのセットアップを行わないようにする
//
// この関数は各コマンドソースの init 関数内で呼び出され、その時点ではまだコマンドの階層化が完了していないため、
// 直接 cmd.CommandPath() を呼び出しただけだと正確なコマンドパスを取得できない。
//
// そのため、この関数は cmd.CommandPath を引数に取り、 SetupWithContext のタイミングで評価することで
// 正確なコマンドパスを取得できるようにしている。
func Ignore(commandPath func() string) {
	ignoreTargets = append(ignoreTargets, commandPath)
}
