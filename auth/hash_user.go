package auth

// HashUser store the hash info of User
type HashUser struct {
	// UserNameHash a 16 bytes hash of username
	UserNameHash string
	// PasswordHash a 16 bytes hash of password
	PasswordHash string
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
