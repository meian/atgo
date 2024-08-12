package cmd

import (
	"os"
	"os/exec"

	"github.com/meian/atgo/files"
	"github.com/meian/atgo/help"
	"github.com/meian/atgo/models/ids"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/usecase"
	"github.com/meian/atgo/workspace"
	"github.com/spf13/cobra"
)

var taskLocalInitFlag struct {
	contestID string
	noOpen    bool
}

var taskLocalInitCmd = &cobra.Command{
	Use:     "local-init [task ID]",
	Aliases: []string{"init"},
	Short:   "Initialize task files on workspace",
	Long: `Initialize task files on workspace.
Create the following file to answer the task for which executed ` + "`{{ .task }}`" + ` at the last.

- main.go (answer file)
- main_test.go (testing for main.go)
- go.mod / go.sum
- task-local.yaml (management file)

Previously created locally files are saved to management path, and if already exists the same task files,
restore from saved it.
If you execute command in vscode, open main.go and main_test.go in vscode after generating files.

Do not edit task-local.yaml manually as it is used when submitting code with ` + "`{{ .submit }}`.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var taskID string
		if len(args) > 0 {
			taskID = args[0]
		}

		p := usecase.TaskLocalInitParam{
			ContestID: ids.ContestID(taskLocalInitFlag.contestID),
			TaskID:    ids.TaskID(taskID),
		}
		res, err := usecase.TaskLocalInit{}.Run(cmd.Context(), p)
		if err != nil {
			return err
		}

		cmd.Println("local workspace is inited.")
		cmd.Println(url.TaskURL(ids.ContestID(res.ContestID), ids.TaskID(res.TaskID)))

		if !taskLocalInitFlag.noOpen && os.Getenv("TERM_PROGRAM") == "vscode" {
			ws := workspace.Dir()
			mainFile := files.MainFile(ws)
			testFile := files.TestFile(ws)
			openCmd := exec.Command("code", mainFile, testFile)
			if err := openCmd.Run(); err != nil {
				return err
			}
		}

		return nil

	},
}

func init() {
	help.Replace(taskLocalInitCmd, []help.Replacer{
		{Key: "task", Command: taskCmd},
		{Key: "submit", Command: submitCmd},
	})
	taskCmd.AddCommand(taskLocalInitCmd)

	taskLocalInitCmd.Flags().StringVarP(&taskLocalInitFlag.contestID, "contest-id", "c", "", "contest ID")
	taskLocalInitCmd.Flags().BoolVar(&taskLocalInitFlag.noOpen, "no-open", false, "do not open go file in editor")
}
