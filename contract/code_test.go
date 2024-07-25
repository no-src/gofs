package contract

import "testing"

func TestCode(t *testing.T) {
	testCases := []struct {
		code Code
		desc string
	}{
		{-99999, UnknownDesc},
		{Unknown, UnknownDesc},
		{Success, SuccessDesc},
		{Fail, FailDesc},
		{Unauthorized, UnauthorizedDesc},
		{NotFound, NotFoundDesc},
		{NoPermission, NoPermissionDesc},
		{ServerError, ServerErrorDesc},
		{AccessDeny, AccessDenyDesc},
		{NotModified, NotModifiedDesc},
		{ChunkNotModified, ChunkNotModifiedDesc},
		{Modified, ModifiedDesc},
		{ChunkModified, ChunkModifiedDesc},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.code.String()
			if actual != tc.desc {
				t.Errorf("code [%d] => expect: %v, but actual: %v \n", tc.code, tc.desc, actual)
			}
		})
	}
}
