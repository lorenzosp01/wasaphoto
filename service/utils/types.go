package utils

import "regexp"

type UserIdentifier struct {
	Id int64 `json:"identifier"`
}

type Username struct {
	Username string `json:"username"`
}

func (u Username) IsValid() bool {
	var valid = false
	if len(u.Username) > 0 && len(u.Username) < 16 {
		valid, _ = regexp.Match(`^[a-zA-Z0-9]+$`, []byte(u.Username))
	}
	return valid
}
