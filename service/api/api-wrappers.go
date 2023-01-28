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

		//todo rimosso il controllo sull'entità del path, metterlo nell'handler
		//var entitiesId [maxParams]int64
		var err error
		var dbErr database.DbError
		params := make(map[string]int64)

		//for i, table := range dbTables {
		//	entitiesId[i], err = strconv.ParseInt(ps[i].Value, 10, 64)
		//	if err != nil {
		//		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400, Message: "Bad request"})
		//		return
		//	}
		//	dbErr = rt.db.EntityExists(entitiesId[i], table)
		//	if dbErr.InternalError != nil {
		//		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		//		return
		//	}
		//	params[ps[i].Key] = entitiesId[i]
		//}

		authorizationHeader := r.Header.Get("Authorization")
		token := utils.GetAuthenticationToken(authorizationHeader)
		if !token.IsValid() {
			rt.LoggerAndHttpErrorSender(w, errors.New("token invalid"), utils.HttpError{StatusCode: 400, Message: "Token is not valid"})
			return
		}

		tokenInt, _ := strconv.ParseInt(token.Value, 10, 64)
		dbErr = rt.db.EntityExists(tokenInt, database.UserTable)
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		}

		params["token"] = tokenInt

		for _, pathParam := range ps {
			params[pathParam.Key], err = strconv.ParseInt(pathParam.Value, 10, 64)
			if err != nil {
				rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400, Message: "Bad request"})
				return
			}
		}

		//todo move query string params in each handler
		//for paramKey, param := range r.URL.Query() {
		//	params[paramKey], err = strconv.ParseInt(param[0], 10, 64)
		//	if err != nil && paramKey != "pattern" {
		//		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 401, Message: "Bad request"})
		//		return
		//	}
		//}

		fn(w, r, params)
	}
}

// Check if user identifier in url path matches the one in the token in the request header
func (rt *_router) authWrap(fn httpRouterHandler) func(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	return func(w http.ResponseWriter, r *http.Request, params map[string]int64) {

		isAuthorized := utils.Authorize(params["token"], params["user_id"])

		if !isAuthorized {
			rt.LoggerAndHttpErrorSender(w, errors.New("user not authorized"), utils.HttpError{StatusCode: 403, Message: "You can't impersonate other users"})
			return
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, params)
	}
}
