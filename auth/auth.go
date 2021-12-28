package auth

import (
	"fmt"
	"github.com/no-src/gofs/contract"
)

// ParseAuthCommandData parse auth command request data
func ParseAuthCommandData(data []byte) (user *HashUser, err error) {
	authCmdLen := len(contract.AuthCommand)
	length := authCmdLen + versionLength + userNameHashLength + PasswordHashLength + expireLength
	if len(data) != length {
		return nil, fmt.Errorf("auth command data is invalid => [%s]", string(data))
	}
	user = &HashUser{
		Version:      data[authCmdLen : authCmdLen+versionLength],
		UserNameHash: string(data[authCmdLen+versionLength : authCmdLen+versionLength+userNameHashLength]),
		PasswordHash: string(data[authCmdLen+versionLength+userNameHashLength : authCmdLen+versionLength+userNameHashLength+PasswordHashLength]),
		Expire:       string(data[authCmdLen+versionLength+userNameHashLength+PasswordHashLength : length]),
	}
	return user, nil
}

// GenerateAuthCommandData generate auth command request data
func GenerateAuthCommandData(user *HashUser) []byte {
	if user == nil {
		return nil
	}
	authData := contract.AuthCommand
	authData = append(authData, user.Version...)
	authData = append(authData, []byte(user.UserNameHash)...)
	authData = append(authData, []byte(user.PasswordHash)...)
	authData = append(authData, []byte(user.Expire)...)
	return authData
}
