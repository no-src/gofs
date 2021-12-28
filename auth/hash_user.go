package auth

import (
	"time"
)

const (
	defaultExpireDuration = time.Second * 15
	userNameHashLength    = 16
	PasswordHashLength    = 16
	expireLength          = 14
	versionLength         = 2
	expireTimeFormat      = "20060102150405"
)

var (
	authVersion = []byte{0, 1}
)

// HashUser store the hash info of User
type HashUser struct {
	// UserNameHash a 16 bytes hash of username
	UserNameHash string
	// PasswordHash a 16 bytes hash of password
	PasswordHash string
	// Expire 14 bytes auth request info expire time of utc, format like "20060102150405"
	Expire string
	// Version 2 bytes of auth api version
	Version []byte
}

// IsExpired auth request info is expired or not
func (h *HashUser) IsExpired() bool {
	if len(h.Expire) != expireLength {
		return true
	}
	expire, err := time.Parse(expireTimeFormat, h.Expire)
	if err != nil {
		return true
	}
	return time.Now().UTC().After(expire)
}

func NewHashUser(userNameHash, passwordHash string) *HashUser {
	return &HashUser{
		UserNameHash: userNameHash,
		PasswordHash: passwordHash,
		Expire:       time.Now().UTC().Add(defaultExpireDuration).Format(expireTimeFormat),
		Version:      authVersion,
	}
}

// ToHashUserList convert User list to HashUser list
func ToHashUserList(users []*User) (hashUsers []*HashUser, err error) {
	if len(users) == 0 {
		return hashUsers, nil
	}
	for _, user := range users {
		hashUser, err := user.ToHashUser()
		if err != nil {
			return nil, err
		}
		hashUsers = append(hashUsers, hashUser)
	}
	return hashUsers, nil
}
