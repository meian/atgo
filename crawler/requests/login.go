package requests

import "net/url"

type Login struct {
	Username  string
	Password  string
	CSRFToken string
	Continue  string
}

func (r Login) URLValues() url.Values {
	return url.Values{
		"username":   {r.Username},
		"password":   {r.Password},
		"csrf_token": {r.CSRFToken},
	}
}
