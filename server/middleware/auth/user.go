package auth

import (
	"encoding/gob"
	"strings"
)

type User struct {
	UserId   int
	UserName string
	Password string
}

// ParseUsers parse users string to User List
// For example: user1|password1,user2|password2
func ParseUsers(userStr string) []*User {
	var users []*User
	if len(userStr) == 0 {
		return users
	}
	all := strings.Split(userStr, ",")
	userCount := 0
	for _, user := range all {
		userInfo := strings.Split(user, "|")
		if len(userInfo) == 2 {
			userName := strings.TrimSpace(userInfo[0])
			password := strings.TrimSpace(userInfo[1])
			if len(userName) > 0 && len(password) > 0 {
				userCount++
				users = append(users, &User{
					UserId:   userCount,
					UserName: userName,
					Password: password,
				})
			}
		}
	}
	return users
}

func init() {
	gob.Register(User{})
}
