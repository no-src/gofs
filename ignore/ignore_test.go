package ignore

import (
	"testing"

	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/util/osutil"
)

const (
	testIgnoreFile = "./testdata/demo.ignore"
)

func TestMatch(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	pi, err := NewPathIgnore(testIgnoreFile, true, logger)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}

	testCases := []struct {
		path   string
		expect bool
	}{
		// for filepath rule
		{"/gofs.exe", true},
		{"/gofs.exe.bak", false},

		{"/debug/", true},
		{"/debug/xx.dll", true},
		{"/debug/subdir", true},
		{"/root/debug/", false},
		{"/root/debug/xx.dll", false},
		{"/root/debug/subdir", false},

		{"/log/gofs1.log", true},
		{"/log/gofs2.log", true},
		{"/log/gofs.log", false},
		{"/log/gofs11.log", false},
		{"/root/log/gofs1.log", false},
		{"/root/log/gofs2.log", false},
		{"/root/log/gofs1.log", false},
		{"/root/log/gofs2.log", false},

		// for regexp rule
		{"/source/.hello.swp", true},
		{"/source/.hello.swp2", false},
		{"bin", true},
		{"/bin", true},
		{"/bin/", true},
		{"bin/", true},
		{"/hello_bin", false},
		{"/source/bin/", true},
		{"/source/bin", true},
		{"/source/bin/hello.txt", true},
		{"/source/bin.log", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			testMatch(t, pi, tc.expect, tc.path)
		})
	}
}

func TestMatch_Windows(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	pi, err := NewPathIgnore(testIgnoreFile, true, logger)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}

	testCases := []struct {
		path   string
		expect bool
	}{
		// for filepath rule
		{"C:\\workspace\\logs\\info.log", true},
		{"C:\\workspace\\logs\\", true},
		{"C:\\workspace\\logs\\2022\\info.log", false},

		{"C:\\workspace\\data\\2022\\my.db", true},
		{"C:\\workspace\\data\\2022\\06", true},
		{"C:\\workspace\\data\\2022\\06\\", false},
		{"C:\\workspace\\data\\my.db", false},
		{"C:\\workspace\\data\\2022\\06\\my.db", false},

		{"C:\\workspace\\doc\\README.MD", true},
		{"C:\\workspace\\doc\\README-CN.MD", true},
		{"C:\\workspace\\doc\\.MD", true},
		{"C:\\workspace\\doc\\*.MD", true},
		{"C:\\workspace\\doc\\README.md", false},
		{"C:\\workspace\\doc\\README-CN.md", false},
		{"C:\\workspace\\doc\\.md", false},
		{"C:\\workspace\\doc\\*.md", false},
	}

	if osutil.IsWindows() {
		for _, tc := range testCases {
			t.Run(tc.path, func(t *testing.T) {
				testMatch(t, pi, tc.expect, tc.path)
			})
		}
	}
}

func TestMatchPath_WithIgnoreDeletedPath_True(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	pi, err := NewPathIgnore(testIgnoreFile, true, logger)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}

	testCases := []struct {
		path   string
		expect bool
	}{
		{"/hello.txt.1643351810.deleted", true},
		{"/hello.txt", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			testMatchPath(t, pi, tc.expect, tc.path)
		})
	}
}

func TestMatchPath_WithIgnoreDeletedPath_False(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	pi, err := NewPathIgnore(testIgnoreFile, false, logger)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}

	testCases := []struct {
		path   string
		expect bool
	}{
		{"/hello.txt.1643351810.deleted", false},
		{"/hello.txt", false},
		{"/source/bin/", true},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			testMatchPath(t, pi, tc.expect, tc.path)
		})
	}
}

func TestParse_ReturnError(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	testCases := []struct {
		expr string
	}{
		{filePathSwitch + "\n*[]"},
		{regexpSwitch + "\n/error**"},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			_, err := parse([]byte(tc.expr), logger)
			if err == nil {
				t.Errorf("parse the rule text should be return error => [%s] error => %v", tc.expr, err)
				return
			}
		})
	}
}

func TestInit_ReturnError(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	c := "./testdata/notfound.ignore"
	_, err := NewPathIgnore(c, true, logger)
	if err == nil {
		t.Errorf("init default ignore component should be return error => %s", c)
		return
	}
}

func TestInit_WithNoConfig(t *testing.T) {
	logger := logger.NewTestLogger()
	defer logger.Close()

	pi, err := NewPathIgnore("", true, logger)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}
	testMatch(t, pi, false, "bin")
}

func testMatch(t *testing.T, pi PathIgnore, expect bool, path string) {
	actual := pi.(*pathIgnore).match(path)
	if actual != expect {
		t.Errorf("[%s] => expect: %v, but actual: %v", path, expect, actual)
	}
}

func testMatchPath(t *testing.T, pi PathIgnore, expect bool, path string) {
	actual := pi.MatchPath(path, "test suit", "test")
	if actual != expect {
		t.Errorf("[%s] => expect: %v, but actual: %v", path, expect, actual)
	}
}
