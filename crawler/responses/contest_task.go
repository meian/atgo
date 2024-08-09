package responses

import (
	"time"

	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
)

type ContestTask struct {
	ContestID ids.ContestID
	Tasks     []ContestTask_Task
}

func (ct ContestTask) ToModel() []models.ContestTask {
	var tasks []models.ContestTask
	for i, t := range ct.Tasks {
		tasks = append(tasks, models.ContestTask{
			ID:        ct.ContestID.ContestTaskID(t.ID),
			ContestID: ct.ContestID,
			TaskID:    t.ID,
			Order:     i + 1,
			Index:     t.Index,
			Task: models.Task{
				ID:        t.ID,
				Title:     t.Title,
				TimeLimit: t.TimeLimit,
				Memory:    t.Memory,
			},
		})
	}
	return tasks
}

type ContestTask_Task struct {
	ID        ids.TaskID
	Title     string
	Index     string
	TimeLimit time.Duration
	Memory    int
}
