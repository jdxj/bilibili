package utils

import (
	"net/http"
)

const (
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36"
)

func StringToCookies(domain, cookie string) []*http.Cookie {
	header := http.Header{}
	header.Set("Cookie", cookie)
	req := &http.Request{Header: header}

	cookies := req.Cookies()
	for _, c := range cookies {
		c.Domain = domain
	}
	return cookies
}

func NewRequestGet(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", UserAgent)
	return req
}
