package cookie

import (
	"net/http"
	"net/url"
	"slices"
)

type JarOption struct {
	IgnorePaths []string
}

type Jar struct {
	jar    http.CookieJar
	option JarOption
}

var _ http.CookieJar = &Jar{}

func NewJar(jar http.CookieJar, opt JarOption) *Jar {
	return &Jar{
		jar:    jar,
		option: opt,
	}

}

// Cookies implements http.CookieJar.
func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	if u == nil {
		return nil
	}
	if slices.Contains(j.option.IgnorePaths, u.Path) {
		return nil
	}
	return j.jar.Cookies(u)
}

// SetCookies implements http.CookieJar.
func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.jar.SetCookies(u, cookies)
}
