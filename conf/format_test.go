package conf

import (
	"path/filepath"
	"testing"
)

func TestFormat_Name(t *testing.T) {
	testCases := []struct {
		format Format
		expect string
	}{
		{JsonFormat, "json"},
		{YamlFormat, "yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			actual := tc.format.Name()
			if actual != tc.expect {
				t.Errorf("test format name error expect:%s, actual:%s", tc.expect, actual)
			}
		})
	}
}

func TestFormat_MatchExt(t *testing.T) {
	testCases := []struct {
		path   string
		format Format
		expect bool
	}{
		{jsonConfigPath, JsonFormat, true},
		{yamlConfigPath, YamlFormat, true},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if !tc.format.MatchExt(filepath.Ext(tc.path)) {
				t.Errorf("match %s confile error => %s", tc.format.Name(), tc.path)
			}
		})
	}
}
