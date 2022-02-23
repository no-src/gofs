package auth

import "testing"

func TestToHashUserList(t *testing.T) {
	// normal convert
	var users []*User
	users = append(users, newUserNoError(t, 1, "root", "toor", DefaultPerm))
	users = append(users, newUserNoError(t, 2, "guest", "guest", DefaultPerm))
	hashUsers, err := ToHashUserList(users)
	if err != nil {
		t.Errorf("convert to hash user list error => %v", err)
		return
	}
	if len(hashUsers) != len(users) {
		t.Errorf("convert to hash user list error")
		return
	}

	// contain nil user
	_, err = ToHashUserList(append(users, nil))
	if err == nil {
		t.Errorf("convert to hash user list should be return error => %v", err)
		return
	}

	// empty user list
	hashUsers, err = ToHashUserList(nil)
	if err != nil {
		t.Errorf("convert an empty user list to hash user list error => %v", err)
		return
	}
	if len(hashUsers) != 0 {
		t.Errorf("convert an empty user list to hash user list should be return an empty hash user list")
		return
	}
}

func TestIsExpired(t *testing.T) {
	hashUser := NewHashUser("698d51a19d8a121c", "bcbe3365e6ac95ea", DefaultPerm)
	if hashUser.IsExpired() {
		t.Errorf("current hash user should not be expired => %s %s", hashUser.UserNameHash, hashUser.Expires)
		return
	}

	hashUser.Expires = "invalid"
	if !hashUser.IsExpired() {
		t.Errorf("current hash user should be expired => %s %s", hashUser.UserNameHash, hashUser.Expires)
		return
	}

	hashUser.Expires = "2006010215040x"
	if !hashUser.IsExpired() {
		t.Errorf("current hash user should be expired => %s %s", hashUser.UserNameHash, hashUser.Expires)
		return
	}
}
