package auth

import "testing"

func TestCheckTo(t *testing.T) {
	testCases := []struct {
		current string
		target  string
		expect  bool
	}{
		{"", "", false},
		{"", "rwx", false},
		{"  ", "", false},
		{"  ", "rwx", false},

		{"r", "rwx", true},
		{"w", "rwx", true},
		{"x", "rwx", true},
		{"rw", "rwx", true},
		{"rwx", "rwx", true},

		{"rwx", "r", false},
		{"rwx", "w", false},
		{"rwx", "x", false},
		{"rwx", "rw", false},
		{"rwx", "rx", false},
		{"rwx", "wx", false},

		{"r", "w", false},
		{"w", "x", false},
		{"x", "r", false},
	}

	for _, tc := range testCases {
		t.Run("["+tc.current+"]=>["+tc.target+"]", func(t *testing.T) {
			cPerm := ToPerm(tc.current)
			tPerm := ToPerm(tc.target)
			actual := cPerm.CheckTo(tPerm)
			if actual != tc.expect {
				t.Errorf("[%s=>%s] => expect: %v, but actual: %v \n", tc.current, tc.target, tc.expect, actual)
			}
		})
	}
}

func TestToPermWithDefault(t *testing.T) {
	testCases := []struct {
		perm   string
		expect string
	}{
		{"r", "r"},
		{"w", "w"},
		{"x", "x"},
		{"", "r"},
		{"a", ""},
		{"abc", ""},
		{"abcd", ""},
		{"rrrr", ""},
	}

	for _, tc := range testCases {
		t.Run("["+tc.perm+"]=>["+tc.expect+"]", func(t *testing.T) {
			cPerm := ToPermWithDefault(tc.perm, ReadPerm)
			actual := cPerm.String()
			if actual != tc.expect {
				t.Errorf("[%s] => expect: %v, but actual: %v \n", tc.perm, tc.expect, actual)
			}
		})
	}
}
