package responses

type Task struct {
	ID        string
	Score     *int
	Samples   []Task_Sample
	CSRFToken string
	LoggedIn  bool
}

type Task_Sample struct {
	Input  string
	Output string
}
