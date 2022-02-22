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

func TestActionInt(t *testing.T) {
	testActionInt(t, 0, UnknownAction)
	testActionInt(t, 1, CreateAction)
	testActionInt(t, 2, WriteAction)
	testActionInt(t, 3, RemoveAction)
	testActionInt(t, 4, RenameAction)
	testActionInt(t, 5, ChmodAction)
	testActionInt(t, 10, Action(10))
	testActionInt(t, UnknownAction.Int(), Action(10).Valid())
}
func testActionInt(t *testing.T, expect int, action Action) {
	actual := action.Int()
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", action, expect, actual)
		t.Fail()
	}
}

func TestActionString(t *testing.T) {
	testActionString(t, "Unknown", UnknownAction)
	testActionString(t, "Create", CreateAction)
	testActionString(t, "Write", WriteAction)
	testActionString(t, "Remove", RemoveAction)
	testActionString(t, "Rename", RenameAction)
	testActionString(t, "Chmod", ChmodAction)
	testActionString(t, "Invalid", Action(10))
}
func testActionString(t *testing.T, expect string, action Action) {
	actual := action.String()
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", action, expect, actual)
		t.Fail()
	}
}
