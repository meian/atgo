package cookie

import (
	"encoding/gob"
	"net/http"
	"net/url"
	"os"
)

var (
	ignoreLoads []func() string
	ignoreSaves []func() string
)

func LoadFrom(url *url.URL, file string, jar http.CookieJar) error {
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

func SaveTo(url *url.URL, file string, jar http.CookieJar) error {
	cookies := jar.Cookies(url)
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

func IgnoreLoad(pathFunc func() string) {
	ignoreLoads = append(ignoreLoads, pathFunc)
}

func ShouldLoad(path string) bool {
	for _, p := range ignoreLoads {
		if p() == path {
			return false
		}
	}
	return true
}

func IgnoreSave(pathFunc func() string) {
	ignoreSaves = append(ignoreSaves, pathFunc)
}

func ShouldSave(path string) bool {
	for _, p := range ignoreSaves {
		if p() == path {
			return false
		}
	}
	return true
}
