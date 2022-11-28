package utils

import (
	"strings"
)

const AuthorizationSchema = "Bearer"

func GetAuthenticationToken(authorizationHeader string) Token {
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 && bearerToken[0] == AuthorizationSchema {
			return Token{bearerToken[1]}
		}
	}

	return Token{}
}

func Authorize(authorizationHeader string, urlId string) HttpError {
	var error HttpError
	token := GetAuthenticationToken(authorizationHeader)

	if token.IsValid() {
		if token.Value == urlId {
			return error
		} else {
			error.StatusCode = 403
			error.Message = "Forbidden"
			return error
		}
	}

	error.StatusCode = 401
	error.Message = "Unauthorized"

	return error
}
