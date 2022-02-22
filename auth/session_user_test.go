package auth

import "testing"

func TestMapperToSessionUserReturnNil(t *testing.T) {
	actual := MapperToSessionUser(nil)
	if actual != nil {
		t.Errorf("MapperToSessionUser should be return nil")
	}
}

func TestMapperToSessionUserReturn(t *testing.T) {
	expect, err := NewUser(1, "root", "toor", "r")
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
}
