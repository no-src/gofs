package fs

import (
	"testing"
)

func TestIsDeleted(t *testing.T) {
	testIsDeleted(t, false, "/test/dir")
	testIsDeleted(t, false, "/test/README.MD")
	testIsDeleted(t, false, "./test/dir")
	testIsDeleted(t, false, "./test/README.MD")

	testIsDeleted(t, false, ".1643351810.deleted")
	testIsDeleted(t, false, "/test/README.MD.164335181.deleted")
	testIsDeleted(t, false, "./test/dir.164335181.deleted")
	testIsDeleted(t, false, "./test/README.MD.164335181.deleted")

	testIsDeleted(t, true, "/test/dir.1643351810.deleted")
	testIsDeleted(t, true, "/test/README.MD.1643351810.deleted")
	testIsDeleted(t, true, "./test/dir.1643351810.deleted")
	testIsDeleted(t, true, "./test/README.MD.1643351810.deleted")

	testIsDeleted(t, true, "/test/dir.16433518101.deleted")
	testIsDeleted(t, true, "/test/README.MD.16433518101.deleted")
	testIsDeleted(t, true, "./test/dir.16433518101.deleted")
	testIsDeleted(t, true, "./test/README.MD.16433518101.deleted")

	testIsDeleted(t, false, ".1643351810.DELETED")
	testIsDeleted(t, true, "/test/dir.1643351810.DELETED")

	testIsDeleted(t, false, "C:\\test\\dir")
	testIsDeleted(t, false, "C:\\test\\README.MD")
	testIsDeleted(t, false, ".\\test\\dir")
	testIsDeleted(t, false, ".\\test\\README.MD")

	testIsDeleted(t, false, "C:\\test\\README.MD.164335181.deleted")
	testIsDeleted(t, false, ".\\test\\dir.164335181.deleted")
	testIsDeleted(t, false, ".\\test\\README.MD.164335181.deleted")

	testIsDeleted(t, true, "C:\\test\\dir.1643351810.deleted")
	testIsDeleted(t, true, "C:\\test\\README.MD.1643351810.deleted")
	testIsDeleted(t, true, ".\\test\\dir.1643351810.deleted")
	testIsDeleted(t, true, ".\\test\\README.MD.1643351810.deleted")

	testIsDeleted(t, true, "C:\\test\\dir.16433518101.deleted")
	testIsDeleted(t, true, "C:\\test\\README.MD.16433518101.deleted")
	testIsDeleted(t, true, ".\\test\\dir.16433518101.deleted")
	testIsDeleted(t, true, ".\\test\\README.MD.16433518101.deleted")

	testIsDeleted(t, false, "C:\\test\\README.MD.164335181.DELETED")
	testIsDeleted(t, true, "C:\\test\\dir.1643351810.DELETED")

}

func testIsDeleted(t *testing.T, expect bool, path string) {
	actual := IsDeleted(path)
	if actual != expect {
		t.Logf("[%s] => expect: %v, but actual: %v \n", path, expect, actual)
		t.Fail()
	}
}
