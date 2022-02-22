package ignore

import "testing"

func TestMatch(t *testing.T) {
	resetDefaultIgnore()
	err := Init("./testdata/demo.ignore", true)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		t.FailNow()
		return
	}
	testMatch(t, true, "/source/.hello.swp")
	testMatch(t, false, "/source/.hello.swp2")

	testMatch(t, true, "bin")
	testMatch(t, true, "/bin")
	testMatch(t, true, "/bin/")
	testMatch(t, true, "bin/")
	testMatch(t, false, "/hello_bin")
	testMatch(t, true, "/source/bin/")
	testMatch(t, true, "/source/bin")
	testMatch(t, true, "/source/bin/hello.txt")
	testMatch(t, false, "/source/bin.log")
}

func testMatch(t *testing.T, expect bool, path string) {
	actual := Match(path)
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", path, expect, actual)
		t.Fail()
	}
}

func TestMatchPath(t *testing.T) {
	resetDefaultIgnore()
	err := Init("./testdata/demo.ignore", true)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		t.FailNow()
		return
	}
	testMatchPath(t, true, "/hello.txt.1643351810.deleted")
	testMatchPath(t, false, "/hello.txt")

	err = Init("./testdata/demo.ignore", false)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		t.FailNow()
		return
	}
	testMatchPath(t, false, "/hello.txt.1643351810.deleted")
	testMatchPath(t, false, "/hello.txt")
	testMatchPath(t, true, "/source/bin/")
}

func testMatchPath(t *testing.T, expect bool, path string) {
	actual := MatchPath(path, "test suit", "test")
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", path, expect, actual)
		t.Fail()
	}
}

func TestParseError(t *testing.T) {
	text := "/error**"
	_, err := parse([]byte(text))
	if err == nil {
		t.Errorf("parse the rule text should be return error => [%s] error => %v", text, err)
		t.FailNow()
		return
	}
}

func TestInitErrorFileNotFound(t *testing.T) {
	resetDefaultIgnore()
	c := "./testdata/notfound.ignore"
	err := Init(c, true)
	if err == nil {
		t.Errorf("init default ignore component should be return error => %s", c)
		t.FailNow()
		return
	}
}

func TestInitWithNoConfig(t *testing.T) {
	resetDefaultIgnore()
	err := Init("", true)
	if err != nil {
		t.Errorf("init default ignore component error => %v", err)
		t.FailNow()
		return
	}
	testMatch(t, false, "bin")
}

func resetDefaultIgnore() {
	defaultIgnore = nil
	matchIgnoreDeletedPath = false
}
