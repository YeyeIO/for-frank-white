package simplejar

import (
	"fmt"
	http "github.com/saucesteals/fhttp"
	"net/url"
)

type Jar struct {
	cookies map[string]*http.Cookie
}

func NewJar() *Jar {
	j := &Jar{}
	j.cookies = map[string]*http.Cookie{}
	return j
}

func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if cookie.Value == "deleted" {
			continue
		}
		j.cookies[cookie.Name] = cookie
	}
}

func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	var cookies []*http.Cookie

	for _, cookie := range j.cookies {
		cookies = append(cookies, cookie)
	}

	return cookies
}

func (j *Jar) Get(name string) (*http.Cookie, error) {
	cookie, ok := j.cookies[name]

	if !ok {
		return nil, fmt.Errorf("simplejar: cookie %q not found", name)
	}

	return cookie, nil
}
