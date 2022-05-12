package httputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// HeaderContentType the Content-Type http header
	HeaderContentType = "Content-Type"
)

var (
	defaultClient    = &http.Client{}
	noRedirectClient = &http.Client{}
	defaultTransport http.RoundTripper
)

var (
	// ErrAppendCertsFromPemFailed attempts to parse a series of PEM encoded certificates failed
	ErrAppendCertsFromPemFailed = errors.New("append certs from pem failed")
)

// HttpGet get http resource
func HttpGet(url string) (resp *http.Response, err error) {
	return defaultClient.Get(url)
}

// HttpGetWithCookie get http resource with cookies
func HttpGetWithCookie(url string, header http.Header, cookies ...*http.Cookie) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	if len(header) > 0 {
		for k, vs := range header {
			for _, v := range vs {
				req.Header.Set(k, v)
			}
		}
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
	req.Header.Set(HeaderContentType, "application/x-www-form-urlencoded")
	return defaultClient.Do(req)
}

// HttpPostFileChunkWithCookie send a post request with form data, a file chunk and cookies
func HttpPostFileChunkWithCookie(url string, fieldName string, fileName string, data url.Values, chunk []byte, cookies ...*http.Cookie) (resp *http.Response, err error) {
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
			req.AddCookie(cookie)
		}
	}
	req.Header.Set(HeaderContentType, w.FormDataContentType())
	return defaultClient.Do(req)
}

// HttpPostWithoutRedirect send a post request with form data and not auto redirect
func HttpPostWithoutRedirect(url string, data url.Values) (resp *http.Response, err error) {
	return noRedirectClient.PostForm(url, data)
}

// Init init default http util
func Init(insecureSkipVerify bool, certFile string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
	if !insecureSkipVerify {
		roots := x509.NewCertPool()
		pemCerts, err := os.ReadFile(certFile)
		if err != nil {
			return err
		}
		if !roots.AppendCertsFromPEM(pemCerts) {
			return ErrAppendCertsFromPemFailed
		}
		tlsConfig.RootCAs = roots
	}
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
		TLSClientConfig:       tlsConfig,
	}
	defaultClient.Transport = defaultTransport
	noRedirectClient.Transport = defaultTransport
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return nil
}
