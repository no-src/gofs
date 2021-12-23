package util

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

var defaultClient = &http.Client{}
var noRedirectClient = &http.Client{}
var defaultTransport http.RoundTripper

// HttpGet get http resource
func HttpGet(url string) (resp *http.Response, err error) {
	return defaultClient.Get(url)
}

// HttpGetWithCookie get http resource with cookies
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

// HttpPost send a post request with form data
func HttpPost(url string, data url.Values) (resp *http.Response, err error) {
	return defaultClient.PostForm(url, data)
}

// HttpPostWithoutRedirect send a post request with form data and not auto redirect
func HttpPostWithoutRedirect(url string, data url.Values) (resp *http.Response, err error) {
	return noRedirectClient.PostForm(url, data)
}

func init() {
	defaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defaultClient.Transport = defaultTransport
	noRedirectClient.Transport = defaultTransport
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
