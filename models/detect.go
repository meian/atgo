package models

import (
	"strings"

	"github.com/meian/atgo/models/ids"
	"github.com/pkg/errors"
)

func DetectContestID(taskID ids.TaskID) (ids.ContestID, error) {
	fs := strings.Split(string(taskID), "_")
	if len(fs) < 2 {
		return "", errors.New("cannot detect contest ID from task ID")
	}
	return ids.ContestID(strings.Join(fs[:len(fs)-1], "_")), nil
}
