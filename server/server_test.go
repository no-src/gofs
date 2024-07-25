package server

import (
	"fmt"
	"testing"
)

func TestGenerateAddr(t *testing.T) {
	testCases := []struct {
		scheme string
		host   string
		port   int
		expect string
	}{
		{"http", "127.0.0.1", 80, "http://127.0.0.1"},
		{"https", "127.0.0.1", 443, "https://127.0.0.1"},
		{"http", "127.0.0.1", 8080, "http://127.0.0.1:8080"},
		{"https", "127.0.0.1", 555, "https://127.0.0.1:555"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s://%s:%d", tc.scheme, tc.host, tc.port), func(t *testing.T) {
			addr := GenerateAddr(tc.scheme, tc.host, tc.port)
			if addr != tc.expect {
				t.Errorf("expect get %s, but actual get %s", tc.expect, addr)
			}
		})
	}
}
