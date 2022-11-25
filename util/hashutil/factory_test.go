package hashutil

import (
	"testing"
)

func TestAllHashAlgorithms(t *testing.T) {
	input := "hello gopher"
	testCases := []struct {
		algorithm string
		input     string
		expect    string
	}{
		{MD5Hash, input, "0e181d273c0a4593230ce8319eaae9b7"},
		{SHA1Hash, input, "07adeacfcdde2146d26119d79cb0028d831e46e9"},
		{SHA256Hash, input, "f7392c4f40eb32d21e6dc087a00049e914f6ce69b76271f46f2b91c53f35166c"},
		{SHA512Hash, input, "d65b9e26ebef94e8984eeb5bba2558183db4a312ca750495e427b288739eebc67ce732b234210638dabe914a9f7202b2facd248ddb753fc649de83574aa2d231"},
		{CRC32Hash, input, "461b91ae"},
		{CRC64Hash, input, "086879fc3731d9ac"},
		{Adler32Hash, input, "1e6804ba"},
		{FNV132Hash, input, "5e358952"},
		{FNV1A32Hash, input, "6836946c"},
		{FNV164Hash, input, "0039da00cd5a26d2"},
		{FNV1A64Hash, input, "d9ced1ff0e72c2cc"},
		{FNV1128Hash, input, "af0bcdfc672247e14fc2ae047bc55b02"},
		{FNV1A128Hash, input, "1a83a6c3193dcc0fbcc90891f7d3df14"},
	}

	for _, tc := range testCases {
		t.Run(tc.algorithm, func(t *testing.T) {
			if err := InitDefaultHash(tc.algorithm); err != nil {
				t.Errorf("calculate hash with [%s] algorithm error, init default hash failed => %v", tc.algorithm, err)
				return
			}
			if actual := Hash([]byte(input)); actual != tc.expect {
				t.Errorf("calculate hash with [%s] algorithm error, expect: %s, but actual: %s", tc.algorithm, tc.expect, actual)
				return
			}

			h, err := NewHash(tc.algorithm)
			if err != nil {
				t.Errorf("calculate hash with [%s] algorithm error, init custom hash failed => %v", tc.algorithm, err)
				return
			}
			if actual := Hash([]byte(input), h); actual != tc.expect {
				t.Errorf("calculate hash with [%s] algorithm error, expect: %s, but actual: %s", tc.algorithm, tc.expect, actual)
				return
			}
		})
	}

	// reset default hash algorithm
	InitDefaultHash(DefaultHash)
}

func TestInitDefaultHash_WithUnsupportedAlgorithm(t *testing.T) {
	testCases := []struct {
		algorithm string
	}{
		{""},
		{"unknown algorithm"},
	}

	for _, tc := range testCases {
		t.Run(tc.algorithm, func(t *testing.T) {
			if err := InitDefaultHash(tc.algorithm); err == nil {
				t.Errorf("calculate hash with [%s] algorithm error, expect to get an error, but get nil", tc.algorithm)
			}
		})
	}

	// reset default hash algorithm
	InitDefaultHash(DefaultHash)
}

func TestNewHash_WithUnsupportedAlgorithm(t *testing.T) {
	testCases := []struct {
		algorithm string
	}{
		{""},
		{"unknown algorithm"},
	}

	for _, tc := range testCases {
		t.Run(tc.algorithm, func(t *testing.T) {
			if _, err := NewHash(tc.algorithm); err == nil {
				t.Errorf("calculate hash with [%s] algorithm error, expect to get an error, but get nil", tc.algorithm)
			}
		})
	}
}
