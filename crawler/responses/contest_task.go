package responses

import (
	"fmt"
	"time"

	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
)

type ContestTask struct {
	ContestID string
	Tasks     []ContestTask_Task
}

func (ct ContestTask) ToModel() []models.ContestTask {
	var tasks []models.ContestTask
	for i, t := range ct.Tasks {
		id := fmt.Sprintf("%s-%s", ct.ContestID, t.ID)
		tasks = append(tasks, models.ContestTask{
			ID:        ids.ContestTaskID(id),
			ContestID: ids.ContestID(ct.ContestID),
			TaskID:    ids.TaskID(t.ID),
			Order:     i + 1,
			Index:     t.Index,
			Task: models.Task{
				ID:        ids.TaskID(t.ID),
				Title:     t.Title,
				TimeLimit: t.TimeLimit,
				Memory:    t.Memory,
			},
		})
	}
	return tasks
}

type ContestTask_Task struct {
	ID        string
	Title     string
	Index     string
	TimeLimit time.Duration
	Memory    int
}
