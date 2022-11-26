package utils

import (
	"strings"
)

const AuthorizationSchema = "Bearer"

func Authorize(authorizationHeader string, urlId string) HttpError {

	var error HttpError
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 && bearerToken[0] == AuthorizationSchema {
			if bearerToken[1] == urlId {
				return error
			} else {
				error = HttpError{StatusCode: 403, Message: "Forbidden"}
				return error
			}
		}
	}

	error = HttpError{StatusCode: 401, Message: "Unauthorized"}
	return error
}
