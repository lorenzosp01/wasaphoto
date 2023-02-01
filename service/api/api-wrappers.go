package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, map[string]int64)

// Parses the request and checks if path is valid
func (rt *_router) wrap(fn httpRouterHandler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var entitiesId [3]int64
		var err error
		var dbErr database.DbError
		params := make(map[string]int64)

		// Check if entities with the given ids exist
		for i, param := range ps {
			entitiesId[i], err = strconv.ParseInt(param.Value, 10, 64)
			if err != nil {
				rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: http.StatusBadRequest, Message: param.Key + " must be in a valid format"})
				return
			}

			doesItExist, dbErr := rt.db.EntityExists(entitiesId[i], database.ParamsNameToTable[param.Key])
			if dbErr.InternalError != nil {
				rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
				return
			} else if !doesItExist {
				rt.LoggerAndHttpErrorSender(w, errors.New("entity not found"), utils.HttpError{StatusCode: http.StatusNotFound, Message: "Entity with " + param.Key + " not found"})
				return
			}

			params[param.Key] = entitiesId[i]
		}

		authorizationHeader := r.Header.Get("Authorization")
		token := utils.GetAuthenticationToken(authorizationHeader)
		if !token.IsValid() {
			rt.LoggerAndHttpErrorSender(w, nil, utils.HttpError{StatusCode: http.StatusBadRequest, Message: "Token is not valid"})
			return
		}

		tokenInt, _ := strconv.ParseInt(token.Value, 10, 64)
		doesUserExists, dbErr := rt.db.EntityExists(tokenInt, database.UserTable)
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		} else if !doesUserExists {
			rt.LoggerAndHttpErrorSender(w, nil, utils.HttpError{StatusCode: http.StatusNotFound, Message: "User with token not found"})
			return
		}

		params["token"] = tokenInt

		for _, pathParam := range ps {
			params[pathParam.Key], err = strconv.ParseInt(pathParam.Value, 10, 64)
			if err != nil {
				rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: http.StatusBadRequest, Message: "Bad request"})
				return
			}
		}

		fn(w, r, params)
	}
}

// Check if user identifier in url path matches the one in the token in the request header
func (rt *_router) authWrap(fn httpRouterHandler) func(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	return func(w http.ResponseWriter, r *http.Request, params map[string]int64) {

		isAuthorized := utils.Authorize(params["token"], params["user_id"])

		if !isAuthorized {
			rt.LoggerAndHttpErrorSender(w, errors.New("user not authorized"), utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't impersonate other users"})
			return
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, params)
	}
}
