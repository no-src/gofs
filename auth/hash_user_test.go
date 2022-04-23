package auth

import (
	"testing"
)

func TestToHashUserList(t *testing.T) {
	testCases := []struct {
		name  string
		users []*User
	}{
		{"normal user list", getTestUserList(t)},
		{"empty user list", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectLen := len(tc.users)
			hashUsers, err := ToHashUserList(tc.users)
			if err != nil {
				t.Errorf("convert to hash user list error => %v", err)
				return
			}
			if len(hashUsers) != expectLen {
				t.Errorf("convert to hash user list error, expect length:%d, actual length:%d", expectLen, len(hashUsers))
			}
		})
	}
}

func TestToHashUserList_ReturnError(t *testing.T) {
	testCases := []struct {
		name  string
		users []*User
	}{
		{"contain nil user", append(getTestUserList(t), nil)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ToHashUserList(tc.users)
			if err == nil {
				t.Errorf("convert to hash user list expect to get an error but get nil")
			}
		})
	}
}

func TestIsExpired(t *testing.T) {
	hashUser := NewHashUser("698d51a19d8a121c", "bcbe3365e6ac95ea", DefaultPerm)
	testCases := []struct {
		name     string
		hashUser *HashUser
		expires  string
		expect   bool
	}{
		{"expires init", hashUser, hashUser.Expires, false},
		{"expires empty", hashUser, "", true},
		{"expires invalid length", hashUser, "invalid", true},
		{"expires invalid value", hashUser, "2006010215040x", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.hashUser.Expires = tc.expires
			if hashUser.IsExpired() != tc.expect {
				t.Errorf("current hash user[%s:%s] IsExpired expect:%v, actual:%v", hashUser.UserNameHash, hashUser.Expires, tc.expect, hashUser.IsExpired())
			}
		})
	}
}
