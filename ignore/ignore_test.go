package ignore

import "testing"

func TestMatch(t *testing.T) {
	resetDefaultIgnore()
	err := Init("./testdata/demo.ignore", true)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}

	testCases := []struct {
		path   string
		expect bool
	}{
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
			testMatch(t, tc.expect, tc.path)
		})
	}
}

func TestMatchPath_WithIgnoreDeletedPath_True(t *testing.T) {
	resetDefaultIgnore()
	err := Init("./testdata/demo.ignore", true)
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
			testMatchPath(t, tc.expect, tc.path)
		})
	}
}

func TestMatchPath_WithIgnoreDeletedPath_False(t *testing.T) {
	resetDefaultIgnore()
	err := Init("./testdata/demo.ignore", false)
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
			testMatchPath(t, tc.expect, tc.path)
		})
	}
}

func TestParse_ReturnError(t *testing.T) {
	text := "/error**"
	_, err := parse([]byte(text))
	if err == nil {
		t.Errorf("parse the rule text should be return error => [%s] error => %v", text, err)
		return
	}
}

func TestInit_ReturnError(t *testing.T) {
	resetDefaultIgnore()
	c := "./testdata/notfound.ignore"
	err := Init(c, true)
	if err == nil {
		t.Errorf("init default ignore component should be return error => %s", c)
		return
	}
}

func TestInit_WithNoConfig(t *testing.T) {
	resetDefaultIgnore()
	err := Init("", true)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		return
	}
	testMatch(t, false, "bin")
}

func testMatch(t *testing.T, expect bool, path string) {
	actual := Match(path)
	if actual != expect {
		t.Errorf("[%s] => expect: %v, but actual: %v", path, expect, actual)
	}
}

func testMatchPath(t *testing.T, expect bool, path string) {
	actual := MatchPath(path, "test suit", "test")
	if actual != expect {
		t.Errorf("[%s] => expect: %v, but actual: %v", path, expect, actual)
	}
}

func resetDefaultIgnore() {
	defaultIgnore = nil
	matchIgnoreDeletedPath = false
}
