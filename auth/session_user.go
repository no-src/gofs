package auth

import "encoding/gob"

type SessionUser struct {
	UserId   int
	UserName string
	Password string
}

func MapperToSessionUser(user *User) *SessionUser {
	if user == nil {
		return nil
	}
	return &SessionUser{
		UserId:   user.UserId(),
		UserName: user.UserName(),
		Password: user.Password(),
	}
}

func init() {
	gob.Register(SessionUser{})
}
