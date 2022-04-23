package auth

import "testing"

func TestMapperToSessionUser(t *testing.T) {
	testCases := []struct {
		userId   int
		userName string
		password string
		perm     string
	}{
		{1, "root", "toor", "rwx"},
		{2, "visitor", "visitor", "r"},
		{3, "worker", "worker", "rw"},
		{4, "no-permission", "123456", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.userName, func(t *testing.T) {
			expect, err := NewUser(tc.userId, tc.userName, tc.password, tc.perm)
			if err != nil {
				t.Errorf("create user error => %v", err)
				return
			}
			if expect == nil {
				t.Errorf("get a nil user")
				return
			}
			actual := MapperToSessionUser(expect)
			if actual == nil || actual.UserId != expect.UserId() || actual.UserName != expect.UserName() || actual.Password != expect.Password() || actual.Perm != expect.Perm() {
				t.Errorf("MapperToSessionUser error, it is not equal")
			}
		})
	}
}

func TestMapperToSessionUser_WithNilUser(t *testing.T) {
	actual := MapperToSessionUser(nil)
	if actual != nil {
		t.Errorf("MapperToSessionUser should be return nil")
	}
}
