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
}

func testIsDeleted(t *testing.T, expect bool, path string) {
	if IsDeleted(path) != expect {
		t.FailNow()
	}
}
