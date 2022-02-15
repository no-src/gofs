package util

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
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

// HttpPostWithCookie send a post request with form data and cookies
func HttpPostWithCookie(url string, data url.Values, cookies ...*http.Cookie) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return defaultClient.Do(req)
}

// HttpPostFileWithCookie send a post request with form data, file and cookies
func HttpPostFileWithCookie(url string, fieldName, fileName string, data url.Values, cookies ...*http.Cookie) (resp *http.Response, err error) {
	f, err := os.Open(fieldName)
	if err != nil {
		return nil, err
	}
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	for k, v := range data {
		for _, item := range v {
			w.WriteField(k, item)
		}
	}

	fw, err := w.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	if len(cookies) > 0 {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return defaultClient.Do(req)
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
