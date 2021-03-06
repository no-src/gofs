package auth

import (
	"errors"
	"time"
)

const (
	userNameHashLength     = 16
	passwordHashLength     = 16
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
	// Perm the user permission
	Perm Perm
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

// NewHashUser create a HashUser instance
func NewHashUser(userNameHash, passwordHash string, perm Perm) *HashUser {
	h := &HashUser{
		UserNameHash: userNameHash,
		PasswordHash: passwordHash,
		Perm:         perm,
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
		if user == nil {
			return nil, errors.New("get a nil user")
		}
		hashUser := user.ToHashUser()
		hashUsers = append(hashUsers, hashUser)
	}
	return hashUsers, nil
}
