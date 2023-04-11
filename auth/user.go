package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/no-src/gofs/util/randutil"
)

// User a login user info
type User struct {
	userId   int
	userName string
	password string
	perm     Perm
}

// String return format user info
func (user *User) String() string {
	return fmt.Sprintf("%s|%s|%s", user.userName, user.password, user.perm)
}

// UserId return user id
func (user *User) UserId() int {
	return user.userId
}

// UserName return username
func (user *User) UserName() string {
	return user.userName
}

// Password return user password
func (user *User) Password() string {
	return user.password
}

// Perm return user permission
func (user *User) Perm() Perm {
	return user.perm
}

// NewUser create a new user
func NewUser(userId int, userName string, password string, perm string) (*User, error) {
	if userId <= 0 {
		return nil, errors.New("userId must greater than zero")
	}
	p := ToPermWithDefault(perm, DefaultPerm)
	if !p.IsValid() {
		return nil, errors.New("user perm must be the composition of 'r' 'w' 'x' or empty")
	}
	user := &User{
		userId:   userId,
		userName: userName,
		password: password,
		perm:     p,
	}
	err := isValidUser(*user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// isValidUser check username and password is valid or not
func isValidUser(user User) error {
	if len(user.UserName()) == 0 {
		return errors.New("userName can't be empty")
	}
	if len(user.Password()) == 0 {
		return errors.New("password can't be empty")
	}
	if strings.ContainsAny(user.UserName(), ",|") {
		return errors.New("userName can't contain ',' or '|' ")
	}
	if strings.ContainsAny(user.Password(), ",|") {
		return errors.New("password can't contain ',' or '|' ")
	}
	if !user.perm.IsValid() {
		return errors.New("user is no permission")
	}
	return nil
}

// ParseUsers parse users string to User List
// For example: user1|password1|rwx,user2|password2|rwx
func ParseUsers(userStr string) (users []*User, err error) {
	if len(userStr) == 0 {
		return users, nil
	}
	all := strings.Split(userStr, ",")
	userCount := 0
	for _, userStr := range all {
		userInfo := strings.Split(userStr, "|")
		fieldLen := len(userInfo)
		if fieldLen >= 2 && fieldLen <= 3 {
			userName := strings.TrimSpace(userInfo[0])
			password := strings.TrimSpace(userInfo[1])
			if len(userName) > 0 && len(password) > 0 {
				userCount++
				perm := ""
				if fieldLen > 2 {
					perm = strings.TrimSpace(userInfo[2])
				}
				user, err := NewUser(userCount, userName, password, perm)
				if err != nil {
					return nil, err
				}
				users = append(users, user)
			}
		} else {
			return nil, fmt.Errorf("invalid user info => [%s]", userStr)
		}
	}
	return users, nil
}

// RandomUser generate some user with random username and password
// count is user count you want
// userLen is the length of random username, max length is 20
// pwdLen is the length of random password, max length is 20
// perm is the default permission of every random user, like 'rwx'
func RandomUser(count, userLen, pwdLen int, perm string) ([]*User, error) {
	var users []*User
	for i := 1; i <= count; i++ {
		user, err := NewUser(i, randutil.RandomString(userLen), randutil.RandomString(pwdLen), perm)
		if err != nil {
			return nil, fmt.Errorf("generate random user error => %s", err.Error())
		}
		users = append(users, user)
	}
	return users, nil
}

// ParseStringUsers parse user list to user string
func ParseStringUsers(users []*User) (userStr string, err error) {
	if len(users) == 0 {
		return userStr, nil
	}
	var userResultList []string
	for _, user := range users {
		if user == nil {
			continue
		}
		err = isValidUser(*user)
		if err != nil {
			return userStr, err
		}
		userResultList = append(userResultList, user.String())
	}
	userStr = strings.Join(userResultList, ",")
	return userStr, nil
}

// GetAnonymousUser get an anonymous user
func GetAnonymousUser() *User {
	return &User{
		userId:   0,
		userName: "anonymous",
		password: "anonymous",
		perm:     FullPerm,
	}
}
