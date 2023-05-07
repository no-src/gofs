package conf

import (
	"errors"
	"os"
	"testing"

	"github.com/no-src/gofs/util/hashutil"
)

const (
	jsonConfigPath = "./example/gofs-remote-client.json"
	yamlConfigPath = "./example/gofs-remote-server.yaml"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{"json configuration", jsonConfigPath},
		{"yaml configuration", yamlConfigPath},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := Config{}
			err := Parse(tc.path, &c)
			if err != nil {
				t.Errorf("parse configuration error => %s", tc.path)
				return
			}
			if c.ChecksumAlgorithm != hashutil.DefaultHash {
				t.Errorf("checksum_algorithm expect to get %s, but actual get %s", hashutil.DefaultHash, c.ChecksumAlgorithm)
			}
		})
	}
}

func TestParse_ReturnError(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		expect error
	}{
		{"empty path", "", errEmptyConfigPath},
		{"unsupported config format", "./parser_test.go", errUnSupportedConfigFormat},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := Config{}
			err := Parse(tc.path, &c)
			if !errors.Is(err, tc.expect) {
				t.Errorf("parse configuration error, expert err:%s, actual err:%s", tc.expect, err)
			}
		})
	}
}

func TestParse_ReturnError_NotExistPath(t *testing.T) {
	c := Config{}
	path := "./not-exist-file.json"
	err := Parse(path, &c)
	if !os.IsNotExist(err) {
		t.Errorf("parse configuration error, expert error not exist, actual err:%s", err)
	}
}

func TestToString(t *testing.T) {
	testCases := []struct {
		name string
		ext  string
		c    Config
	}{
		{"json configuration", ".json", Config{}},
		{"yaml configuration", ".yaml", Config{}},
		{"yml configuration", ".yml", Config{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ToString(tc.ext, tc.c)
			if err != nil {
				t.Errorf("convert configuration to string error => %s", tc.name)
			}
		})
	}
}

func TestToString_ReturnError(t *testing.T) {
	testCases := []struct {
		name string
		ext  string
		c    Config
	}{
		{"empty ext", "", Config{}},
		{"invalid ext", ".xyz", Config{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ToString(tc.ext, tc.c)
			expect := errUnSupportedConfigFormat
			if !errors.Is(err, expect) {
				t.Errorf("convert configuration to string error, expert err:%s, actual err:%s", expect, err)
			}
		})
	}
}
