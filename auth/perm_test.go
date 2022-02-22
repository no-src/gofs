package auth

import "testing"

func TestCheckTo(t *testing.T) {
	testCheckTo(t, false, "", "")
	testCheckTo(t, false, "", "rwx")
	testCheckTo(t, false, "  ", "")
	testCheckTo(t, false, "  ", "rwx")
	testCheckTo(t, true, "r", "rwx")
	testCheckTo(t, true, "w", "rwx")
	testCheckTo(t, true, "x", "rwx")
	testCheckTo(t, true, "rw", "rwx")
	testCheckTo(t, true, "rwx", "rwx")

	testCheckTo(t, false, "rwx", "r")
	testCheckTo(t, false, "rwx", "w")
	testCheckTo(t, false, "rwx", "x")
	testCheckTo(t, false, "rwx", "rw")
	testCheckTo(t, false, "rwx", "rx")
	testCheckTo(t, false, "rwx", "wx")

	testCheckTo(t, false, "r", "w")
	testCheckTo(t, false, "w", "x")
	testCheckTo(t, false, "x", "r")
}

func testCheckTo(t *testing.T, expect bool, current, target string) {
	cPerm := ToPerm(current)
	tPerm := ToPerm(target)
	actual := cPerm.CheckTo(tPerm)
	if actual != expect {
		t.Logf("[%s=>%s] => expect: %v, but actual: %v \n", current, target, expect, actual)
		t.Fail()
	}
}

func TestToPermWithDefault(t *testing.T) {
	testToPermWithDefault(t, "r", "r")
	testToPermWithDefault(t, "w", "w")
	testToPermWithDefault(t, "x", "x")
	testToPermWithDefault(t, "r", "")
	testToPermWithDefault(t, "", "a")
	testToPermWithDefault(t, "", "abc")
	testToPermWithDefault(t, "", "abcd")
	testToPermWithDefault(t, "", "rrrr")
}

func testToPermWithDefault(t *testing.T, expect string, perm string) {
	cPerm := ToPermWithDefault(perm, ReadPerm)
	actual := cPerm.String()
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", perm, expect, actual)
		t.Fail()
	}
}
