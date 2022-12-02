package utils

import (
	"strconv"
)

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
