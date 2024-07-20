package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/pkg/errors"
)

func Read(ctx context.Context, file string) (string, string, error) {
	logger := logs.FromContext(ctx).With("file", file)
	if !io.FileExists(file) {
		logger.Error("No store credential")
		return "", "", errors.New("No store credential")
	}
	cr, err := os.ReadFile(file)
	if err != nil {
		logger.Error(err.Error())
		return "", "", errors.New("failed to open credential")
	}
	parts := strings.Split(string(cr), "\n")
	if len(parts) < 2 {
		return "", "", errors.New("invalid credential format")
	}
	username := parts[0]
	ep := parts[1]
	password, err := decrypt(ep)
	if err != nil {
		logger.Error(err.Error())
		return "", "", errors.New("failed to read password")
	}
	return username, password, nil
}

func Write(ctx context.Context, file, user, password string) error {
	logger := logs.FromContext(ctx).With("file", file)
	fw, err := os.Create(file)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to create credential")
	}
	defer fw.Close()
	bw := bufio.NewWriter(fw)
	if _, err := fmt.Fprintln(bw, user); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to write username")
	}
	if _, err := fmt.Fprintln(bw, encrypt(password)); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to write password")
	}
	if err := bw.Flush(); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to flush")
	}
	return nil
}
