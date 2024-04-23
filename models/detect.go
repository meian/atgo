package models

import (
	"strings"

	"github.com/pkg/errors"
)

func DetectContestID(taskID string) (string, error) {
	fs := strings.Split(taskID, "_")
	if len(fs) < 2 {
		return "", errors.New("cannot detect contest ID from task ID")
	}
	return strings.Join(fs[:len(fs)-1], "_"), nil
}
