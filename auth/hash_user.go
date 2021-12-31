package auth

import (
	"time"
)

const (
	userNameHashLength     = 16
	PasswordHashLength     = 16
	versionLength          = 2
	expiresLength          = 14
	expiresFormat          = "20060102150405"
	defaultExpiresDuration = time.Second * 15
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
	// Expires 14 bytes auth request info expires of utc, format like "20060102150405"
	Expires string
	// Version 2 bytes of auth api version
	Version []byte
}

// IsExpired auth request info is expired or not
func (h *HashUser) IsExpired() bool {
	if len(h.Expires) != expiresLength {
		return true
	}
	expires, err := time.Parse(expiresFormat, h.Expires)
	if err != nil {
		return true
	}
	return time.Now().UTC().After(expires)
}

// RefreshExpires refresh expires with current utc time
func (h *HashUser) RefreshExpires() string {
	h.Expires = time.Now().UTC().Add(defaultExpiresDuration).Format(expiresFormat)
	return h.Expires
}

func NewHashUser(userNameHash, passwordHash string) *HashUser {
	h := &HashUser{
		UserNameHash: userNameHash,
		PasswordHash: passwordHash,
		Version:      authVersion,
	}
	h.RefreshExpires()
	return h
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
