package util

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

var defaultClient = &http.Client{}
var noRedirectClient = &http.Client{}

func HttpGet(url string) (resp *http.Response, err error) {
	return defaultClient.Get(url)
}

func HttpGetWithCookie(url string, cookies ...*http.Cookie) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return defaultClient.Do(req)
}

func HttpPost(url string, data url.Values) (resp *http.Response, err error) {
	return defaultClient.PostForm(url, data)
}

func HttpPostWithoutRedirect(url string, data url.Values) (resp *http.Response, err error) {
	return noRedirectClient.PostForm(url, data)
}

func init() {
	defaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	noRedirectClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
