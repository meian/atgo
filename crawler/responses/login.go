package responses

type Login struct {
	LoggedIn  bool
	CSRFToken string
}
