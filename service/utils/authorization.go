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

func Authorize(token Token, urlId string) bool {
	return token.Value == urlId
}
