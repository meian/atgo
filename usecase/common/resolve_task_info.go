package common

import (
	"context"
	"errors"

	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/workspace"
)

func ResolveTaskInfo(ctx context.Context, contestID ids.ContestID, taskID ids.TaskID) (*models.TaskInfo, bool, error) {
	logger := logs.FromContext(ctx)
	if len(taskID) > 0 {
		logger = logger.With("taskID", taskID)
		if len(contestID) == 0 {
			cID, err := models.DetectContestID(ids.TaskID(taskID))
			if err != nil {
				logger.Error(err.Error())
				return nil, false, errors.New("failed to detect contest ID from task ID")
			}
			contestID = cID
		}
	}

	file, ok := workspace.TaskInfoFile()
	logger = logger.With("info file", file)
	if !ok {
		if len(contestID) == 0 {
			return nil, false, errors.New("not exists task info file")
		}
		return &models.TaskInfo{
			ContestID: contestID,
			TaskID:    ids.TaskID(taskID),
		}, true, nil
	}
	var info models.TaskInfo
	if err := info.ReadFile(file); err != nil {
		logger.Error(err.Error())
		return nil, false, errors.New("cannot read task info file")
	}
	if len(info.ContestID) == 0 {
		return nil, false, errors.New("cannot get contest ID from task info file")
	}
	mustSave := false
	if len(contestID) > 0 {
		if info.ContestID != contestID {
			info.ContestID = contestID
			info.TaskID = ids.TaskID(taskID)
			mustSave = true
		}
	}
	if len(taskID) > 0 {
		if info.TaskID != ids.TaskID(taskID) {
			info.TaskID = ids.TaskID(taskID)
			mustSave = true
		}
	}

	return &info, mustSave, nil
}
