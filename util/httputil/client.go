package httputil

import (
	"bytes"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/quic-go/quic-go/http3"
)

type httpClient struct {
	defaultClient    *http.Client
	noRedirectClient *http.Client
}

// NewHttpClient create a http client
func NewHttpClient(insecureSkipVerify bool, certFile string, enableHTTP3 bool) (HttpClient, error) {
	c := &httpClient{
		defaultClient:    &http.Client{},
		noRedirectClient: &http.Client{},
	}
	tlsConfig, err := NewTLSConfig(insecureSkipVerify, certFile)
	if err != nil {
		return nil, err
	}

	var rt http.RoundTripper
	if enableHTTP3 {
		rt = &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
		}
	} else {
		rt = &http.Transport{
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
			TLSClientConfig:       tlsConfig,
		}
	}

	c.defaultClient.Transport = rt
	c.noRedirectClient.Transport = rt
	c.noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return c, nil
}

func (c *httpClient) HttpGet(url string) (resp *http.Response, err error) {
	return c.defaultClient.Get(url)
}

func (c *httpClient) HttpGetWithCookie(url string, header http.Header, cookies ...*http.Cookie) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		if cookie != nil {
			req.AddCookie(cookie)
		}
	}

	if len(header) > 0 {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Set(k, v)
			}
		}
	}
	return c.defaultClient.Do(req)
}

func (c *httpClient) HttpPost(url string, data url.Values) (resp *http.Response, err error) {
	return c.defaultClient.PostForm(url, data)
}

func (c *httpClient) HttpPostWithCookie(url string, data url.Values, cookies ...*http.Cookie) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		if cookie != nil {
			req.AddCookie(cookie)
		}
	}
	req.Header.Set(HeaderContentType, "application/x-www-form-urlencoded")
	return c.defaultClient.Do(req)
}

func (c *httpClient) HttpPostFileChunkWithCookie(url string, fieldName string, fileName string, data url.Values, chunk []byte, cookies ...*http.Cookie) (resp *http.Response, err error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	for k, v := range data {
		for _, item := range v {
			w.WriteField(k, item)
		}
	}

	fw, err := w.CreateFormFile(fieldName, filepath.Base(fileName))
	if err != nil {
		return nil, err
	}

	if len(chunk) > 0 {
		if _, err = fw.Write(chunk); err != nil {
			return nil, err
		}
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	if len(cookies) > 0 {
		for _, cookie := range cookies {
			if cookie != nil {
				req.AddCookie(cookie)
			}
		}
	}
	req.Header.Set(HeaderContentType, w.FormDataContentType())
	return c.defaultClient.Do(req)
}

func (c *httpClient) HttpPostWithoutRedirect(url string, data url.Values) (resp *http.Response, err error) {
	return c.noRedirectClient.PostForm(url, data)
}
