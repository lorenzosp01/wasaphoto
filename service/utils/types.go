package utils

import (
	"regexp"
	"strconv"
)

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

type Token struct {
	Value string
}

func (t Token) IsValid() bool {
	if len(t.Value) > 0 {
		_, err := strconv.ParseInt(t.Value, 10, 64)
		return err == nil
	}
	return false
}

type HttpError struct {
	StatusCode int
	Message    string
}
