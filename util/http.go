package util

import (
	"crypto/tls"
	"net/http"
)

var defaultClient = &http.Client{}

func HttpGet(url string) (resp *http.Response, err error) {
	return defaultClient.Get(url)
}

func init() {
	defaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}
