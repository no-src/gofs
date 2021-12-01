package auth

import (
	"fmt"
	"github.com/no-src/gofs/contract"
)

// ParseAuthCommandData parse auth command request data
func ParseAuthCommandData(data []byte) (user *HashUser, err error) {
	authCmdLen := len(contract.AuthCommand)
	length := authCmdLen + 16 + 16
	if len(data) != length {
		return nil, fmt.Errorf("auth command data is invalid => [%s]", string(data))
	}
	user = &HashUser{
		UserNameHash: string(data[authCmdLen : authCmdLen+16]),
		PasswordHash: string(data[authCmdLen+16 : authCmdLen+32]),
	}
	return user, nil
}

// GenerateAuthCommandData generate auth command request data
func GenerateAuthCommandData(userNameHash, passwordHash string) []byte {
	authData := contract.AuthCommand
	authData = append(authData, []byte(userNameHash)...)
	authData = append(authData, []byte(passwordHash)...)
	return authData
}
