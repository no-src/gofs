package action

import "testing"

func TestParseActionFromString(t *testing.T) {
	testParseActionFromString(t, CreateAction, "1")
	testParseActionFromString(t, WriteAction, "2")
	testParseActionFromString(t, RemoveAction, "3")
	testParseActionFromString(t, RenameAction, "4")
	testParseActionFromString(t, ChmodAction, "5")
	testParseActionFromString(t, UnknownAction, "99999")
	testParseActionFromString(t, UnknownAction, "xyz")
	testParseActionFromString(t, UnknownAction, "0")
	testParseActionFromString(t, UnknownAction, "-1")
}

func testParseActionFromString(t *testing.T, expect Action, action string) {
	actual := ParseActionFromString(action)
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", action, expect, actual)
		t.Fail()
	}
}
