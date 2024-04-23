package requests

import "net/url"

type Home struct{}

func (r Home) URLValues() url.Values {
	return url.Values{}
}
