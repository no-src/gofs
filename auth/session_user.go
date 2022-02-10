package auth

import "encoding/gob"

// SessionUser the login user info that is stored in session
type SessionUser struct {
	UserId   int
	UserName string
	Password string
	Perm     Perm
}

// MapperToSessionUser convert User to SessionUser
func MapperToSessionUser(user *User) *SessionUser {
	if user == nil {
		return nil
	}
	return &SessionUser{
		UserId:   user.UserId(),
		UserName: user.UserName(),
		Password: user.Password(),
		Perm:     user.Perm(),
	}
}

func init() {
	gob.Register(SessionUser{})
}
