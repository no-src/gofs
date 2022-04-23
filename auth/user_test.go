package auth

import (
	"testing"
)

func TestRandomUser(t *testing.T) {
	testCases := []struct {
		name    string
		count   int
		userLen int
		pwdLen  int
		perm    string
	}{
		{"user count 0", 0, 8, 8, DefaultPerm},
		{"user count 1", 1, 8, 8, DefaultPerm},
		{"user count 5", 5, 8, 8, DefaultPerm},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			users, err := RandomUser(tc.count, tc.userLen, tc.pwdLen, tc.perm)
			if err != nil {
				t.Errorf("generate random user error => %v", err)
				return
			}
			if len(users) != tc.count {
				t.Errorf("generate random user count expect:%d actual:%d", tc.count, len(users))
			}
		})
	}
}

func TestRandomUser_ReturnError(t *testing.T) {
	testCases := []struct {
		name    string
		count   int
		userLen int
		pwdLen  int
		perm    string
	}{
		{"invalid permission", 5, 8, 8, "abc"},
		{"invalid length of username", 5, 0, 8, DefaultPerm},
		{"invalid length of password", 5, 8, 0, DefaultPerm},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := RandomUser(tc.count, tc.userLen, tc.pwdLen, tc.perm)
			if err == nil {
				t.Errorf("generate random user expect to get an error but get nil")
				return
			}
		})
	}
}

func TestParseStringUsers(t *testing.T) {
	testCases := []struct {
		name  string
		users []*User
	}{
		{"empty user list", []*User{}},
		{"test user root", append([]*User{}, newUserNoError(t, 1, "root", "toor", FullPerm))},
		{"test user guest", append([]*User{}, newUserNoError(t, 2, "guest", "guest", DefaultPerm))},
		{"append nil user", append([]*User{}, nil)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userStr, err := ParseStringUsers(tc.users)
			if err != nil {
				t.Errorf("ParseStringUsers error => %v", err)
				return
			}
			if len(tc.users) == 0 && len(userStr) != 0 {
				t.Errorf("ParseStringUsers parse empty user list error, expect to get an empty userstr, but get %s", userStr)
			}
		})
	}
}

func TestParseStringUsers_ReturnError(t *testing.T) {
	testCases := []struct {
		name  string
		users []*User
	}{
		{"empty username", append(getTestUserList(t), createUserNoValidate(3, "", "pwd", DefaultPerm))},
		{"empty password", append(getTestUserList(t), createUserNoValidate(3, "user", "", DefaultPerm))},
		{"invalid username", append(getTestUserList(t), createUserNoValidate(3, "user|,", "pwd", DefaultPerm))},
		{"invalid password", append(getTestUserList(t), createUserNoValidate(3, "user", "pwd|,", DefaultPerm))},
		{"invalid permission", append(getTestUserList(t), createUserNoValidate(3, "user", "pwd", "abc"))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseStringUsers(tc.users)
			if err == nil {
				t.Errorf("ParseStringUsers expect to get an error but get nil")
			}
		})
	}
}

func TestParseUsers(t *testing.T) {
	testCases := []struct {
		userStr string
		expect  int
	}{
		{"user1|password1|rwx,user2|password2|rwx", 2},
		{"", 0},
	}

	for _, tc := range testCases {
		t.Run("["+tc.userStr+"]", func(t *testing.T) {
			users, err := ParseUsers(tc.userStr)
			if err != nil {
				t.Errorf("parse user error => %v", err)
				return
			}
			if len(users) != tc.expect {
				t.Errorf("parse user expect get %d users, actual:%d ", tc.expect, len(users))
			}
		})
	}
}

func TestParseUsers_ReturnError(t *testing.T) {
	testCases := []struct {
		name    string
		userStr string
	}{
		{"invalid permission", "user1|password1|abc"},
		{"invalid field", "user1|password1|abc|ext"},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"=>"+tc.name, func(t *testing.T) {
			_, err := ParseUsers(tc.userStr)
			if err == nil {
				t.Errorf("ParseUsers expect to get an error but get nil")
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	testCases := []struct {
		userId   int
		userName string
		password string
		perm     string
	}{
		{1, "root", "toor", DefaultPerm},
		{2, "admin", "admin", FullPerm},
	}

	for _, tc := range testCases {
		t.Run(tc.userName, func(t *testing.T) {
			_, err := NewUser(tc.userId, tc.userName, tc.password, tc.perm)
			if err != nil {
				t.Errorf("create user error => %v", err)
			}
		})
	}
}

func TestNewUser_ReturnError(t *testing.T) {
	testCases := []struct {
		name     string
		userId   int
		userName string
		password string
		perm     string
	}{
		{"invalid userid", -1, "root", "toor", DefaultPerm},
		{"empty username", 1, "", "toor", DefaultPerm},
		{"empty password", 1, "root", "", DefaultPerm},
		{"invalid username", 1, "root|,", "toor", DefaultPerm},
		{"invalid password", 1, "root", "toor|,", DefaultPerm},
		{"invalid permission", 1, "root", "toor", "abc"},
	}

	for _, tc := range testCases {
		t.Run(tc.userName, func(t *testing.T) {
			_, err := NewUser(tc.userId, tc.userName, tc.password, tc.perm)
			if err == nil {
				t.Errorf("NewUser expect to get an error but get nil")
			}
		})
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

func createUserNoValidate(userId int, userName string, password string, perm Perm) *User {
	user := &User{
		userId:   userId,
		userName: userName,
		password: password,
		perm:     perm,
	}
	return user
}

func getTestUserList(t *testing.T) []*User {
	var users []*User
	users = append(users, newUserNoError(t, 1, "root", "toor", DefaultPerm))
	users = append(users, newUserNoError(t, 2, "guest", "guest", DefaultPerm))
	return users
}
