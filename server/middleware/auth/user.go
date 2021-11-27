package auth

import "encoding/gob"

type User struct {
	UserId   int
	UserName string
	Password string
}

func init() {
	gob.Register(User{})
}
