package auth

import (
	"bufio"
	"context"
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
	fr, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
		return "", "", errors.New("failed to open credential")
	}
	defer fr.Close()
	br := bufio.NewReader(fr)
	username, err := br.ReadString('\n')
	if err != nil {
		logger.Error(err.Error())
		return "", "", errors.New("failed to read username")
	}
	username = strings.TrimSpace(username)
	ep, err := br.ReadString('\n')
	if err != nil {
		logger.Error(err.Error())
		return "", "", errors.New("failed to read password")
	}
	ep = strings.TrimSpace(ep)
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
	_, err = bw.WriteString(user + "\n")
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to write username")
	}
	ep, err := encrypt(password)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to write password")
	}
	_, err = bw.WriteString(ep + "\n")
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to write password")
	}
	if err := bw.Flush(); err != nil {
		logger.Error(err.Error())
		return errors.New("failed to flush")
	}
	return nil
}
