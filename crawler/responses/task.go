package responses

import "github.com/meian/atgo/models/ids"

type Task struct {
	ID        ids.TaskID
	Score     *int
	Samples   []Task_Sample
	CSRFToken string
	LoggedIn  bool
}

type Task_Sample struct {
	Input  string
	Output string
}
