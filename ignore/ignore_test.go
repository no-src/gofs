package ignore

import "testing"

func TestMatch(t *testing.T) {
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
