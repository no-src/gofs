package auth

import "testing"

func TestRandomUser(t *testing.T) {
	userCount := 5
	users, err := RandomUser(userCount, 8, 8, DefaultPerm)
	if err != nil {
		t.Errorf("generate random user error => %v", err)
		return
	}
	if len(users) != userCount {
		t.Errorf("generate random user count expect:%d actual:%d", userCount, len(users))
		return
	}
}

func TestRandomUserError(t *testing.T) {
	userCount := 5
	_, err := RandomUser(userCount, 8, 8, "abc")
	if err == nil {
		t.Errorf("generate random user should be return error")
		return
	}
}

func TestParseStringUsers(t *testing.T) {
	var users []*User
	users = append(users, newUserNoError(t, 1, "root", "toor", DefaultPerm))
	users = append(users, newUserNoError(t, 2, "guest", "guest", DefaultPerm))
	users = append(users, nil)
	_, err := ParseStringUsers(users)
	if err != nil {
		t.Errorf("ParseStringUsers return error => %v", err)
		return
	}

	// empty username
	u := getUser(t)
	u.userName = ""
	_, err = ParseStringUsers(append(users, u))
	if err == nil {
		t.Errorf("ParseStringUsers should be return error because of the empty username")
		return
	}

	// empty password
	u = getUser(t)
	u.password = ""
	_, err = ParseStringUsers(append(users, u))
	if err == nil {
		t.Errorf("ParseStringUsers should be return error because of the empty password")
		return
	}

	// invalid username
	u = getUser(t)
	u.userName = "user|,"
	_, err = ParseStringUsers(append(users, u))
	if err == nil {
		t.Errorf("ParseStringUsers should be return error because of the invalid username")
		return
	}

	// invalid password
	u = getUser(t)
	u.password = "pwd|,"
	_, err = ParseStringUsers(append(users, u))
	if err == nil {
		t.Errorf("ParseStringUsers should be return error because of the invalid password")
		return
	}

	// invalid permission
	u = getUser(t)
	u.perm = "abc"
	_, err = ParseStringUsers(append(users, u))
	if err == nil {
		t.Errorf("ParseStringUsers should be return error because of the invalid permission")
		return
	}

	userStr, err := ParseStringUsers(nil)
	if err != nil || len(userStr) != 0 {
		t.Errorf("ParseStringUsers parse empty user list error")
		return
	}
}

func getUser(t *testing.T) *User {
	u := newUserNoError(t, 3, "user", "pwd", DefaultPerm)
	return u
}

func TestParseUsers(t *testing.T) {
	users, err := ParseUsers("user1|password1|rwx,user2|password2|rwx")
	if err != nil {
		t.Errorf("parse user error => %v", err)
		return
	}
	expect := 2
	if len(users) != expect {
		t.Errorf("parse user expect get %d users, actual:%d ", expect, len(users))
		return
	}

	_, err = ParseUsers("user1|password1|abc")
	if err == nil {
		t.Errorf("parse user should be return error because of the invalid permission")
		return
	}

	users, err = ParseUsers("")
	if err != nil || len(users) != 0 {
		t.Errorf("parse empty user string should be return nil error and empty users")
		return
	}

	_, err = ParseUsers("user1|password1|abc|ext")
	if err == nil {
		t.Errorf("parse user should be return error because of the invalid field")
		return
	}
}

func TestNewUser(t *testing.T) {

	newUserNoError(t, 1, "root", "toor", DefaultPerm)

	_, err := NewUser(-1, "root", "toor", DefaultPerm)
	if err == nil {
		t.Errorf("create a user should be return error because the invalid userid")
		return
	}

	_, err = NewUser(1, "", "toor", DefaultPerm)
	if err == nil {
		t.Errorf("create a user should be return error because the empty username")
		return
	}
}

func newUserNoError(t *testing.T, userId int, userName string, password string, perm string) *User {
	u, err := NewUser(userId, userName, password, perm)
	if err != nil {
		t.Errorf("create user error => %v", err)
		return nil
	}
	return u
}
