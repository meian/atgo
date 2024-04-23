package cookiestore

import (
	"encoding/gob"
	"net/http"
	"net/url"
	"os"
)

func Load(url *url.URL, file string, jar http.CookieJar) error {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()
	var cookies []*http.Cookie
	if err := gob.NewDecoder(f).Decode(&cookies); err != nil {
		return err
	}
	jar.SetCookies(url, cookies)
	return nil
}

func Save(url *url.URL, file string, jar http.CookieJar) error {
	var cookies []*http.Cookie
	for _, cookie := range jar.Cookies(url) {
		if cookie.MaxAge > 0 {
			cookies = append(cookies, cookie)
		}
	}
	if len(cookies) == 0 {
		return nil
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := gob.NewEncoder(f).Encode(cookies); err != nil {
		return err
	}
	return nil
}
