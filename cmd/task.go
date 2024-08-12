package cmd

import (
	"github.com/meian/atgo/io"
	"github.com/meian/atgo/logs"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/tmpl"
	"github.com/meian/atgo/usecase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var taskFlag struct {
	contestID   string
	showSamples bool
}

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task [task ID]",
	Short: "Display task info",
	Long:  `Display task and associated samples.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var taskID string
		if len(args) > 0 {
			taskID = args[0]
		}

		p := usecase.TaskParam{
			TaskID:      ids.TaskID(taskID),
			ContestID:   ids.ContestID(taskFlag.contestID),
			ShowSamples: taskFlag.showSamples,
		}
		res, err := usecase.Task{}.Run(ctx, p)
		if err != nil {
			return err
		}

		tname := "task"
		t := tmpl.CmdTemplate(tname)
		w := io.OutFromContext(ctx)
		model := toTaskOutputModel(res)
		if err := t.Execute(w, model); err != nil {
			logs.FromContext(ctx).With("template name", tname).Error(err.Error())
			return errors.New("failed to execute template")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.Flags().StringVarP(&taskFlag.contestID, "contest-id", "c", "", "contest ID")
	taskCmd.Flags().BoolVarP(&taskFlag.showSamples, "show-samples", "s", false, "show task samples")
}

func toTaskOutputModel(res *usecase.TaskResult) TaskOutputModel {
	return TaskOutputModel{
		Contest: res.Contest,
		Index:   res.ContestTask.Index,
		Task:    res.Task,
	}
}

type TaskOutputModel struct {
	Contest models.Contest
	Index   string
	Task    models.Task
}
