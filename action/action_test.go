package action

import "testing"

func TestParseActionFromString(t *testing.T) {
	testCases := []struct {
		action string
		expect Action
	}{
		{"1", CreateAction},
		{"2", WriteAction},
		{"3", RemoveAction},
		{"4", RenameAction},
		{"5", ChmodAction},
		{"6", SymlinkAction},
		{"99999", UnknownAction},
		{"xyz", UnknownAction},
		{"0", UnknownAction},
		{"-1", UnknownAction},
	}

	for _, tc := range testCases {
		t.Run(tc.action, func(t *testing.T) {
			actual := ParseActionFromString(tc.action)
			if actual != tc.expect {
				t.Errorf("[%s] => expect: %v, but actual: %v", tc.action, tc.expect, actual)
			}
		})
	}
}

func TestAction_Int(t *testing.T) {
	testCases := []struct {
		name   string
		action Action
		expect int
	}{
		{"UnknownAction", UnknownAction, 0},
		{"CreateAction", CreateAction, 1},
		{"WriteAction", WriteAction, 2},
		{"RemoveAction", RemoveAction, 3},
		{"RenameAction", RenameAction, 4},
		{"ChmodAction", ChmodAction, 5},
		{"SymlinkAction", SymlinkAction, 6},
		{"Action(10)", Action(10), 10},
		{"Action(10).Valid()", Action(10).Valid(), UnknownAction.Int()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.action.Int()
			if actual != tc.expect {
				t.Errorf("[%s] => expect: %v, but actual: %v", tc.action, tc.expect, actual)
			}
		})
	}
}

func TestAction_String(t *testing.T) {
	testCases := []struct {
		name   string
		action Action
		expect string
	}{
		{"UnknownAction", UnknownAction, "Unknown"},
		{"CreateAction", CreateAction, "Create"},
		{"WriteAction", WriteAction, "Write"},
		{"RemoveAction", RemoveAction, "Remove"},
		{"RenameAction", RenameAction, "Rename"},
		{"ChmodAction", ChmodAction, "Chmod"},
		{"SymlinkAction", SymlinkAction, "Symlink"},
		{"Action(10)", Action(10), "Invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.action.String()
			if actual != tc.expect {
				t.Errorf("[%s] => expect: %v, but actual: %v", tc.action, tc.expect, actual)
			}
		})
	}
}
