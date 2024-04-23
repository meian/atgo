package url

import "net/url"

type Valuer interface {
	URLValues() url.Values
}
