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
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, utils.Token)

// Parses the request and checks if path in the id are valid and entity with that exists
//todo passare a fn ad una struct con i parametri nel path e del querystring
func (rt *_router) wrap(fn httpRouterHandler, dbTables []string, filterParams bool) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var entitiesId [3]int64
		var err error
		var dbErr database.DbError
		for i, table := range dbTables {
			entitiesId[i], err = strconv.ParseInt(ps[i].Value, 10, 64)
			if err != nil {
				rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
				return
			}

			dbErr = rt.db.EntityExists(entitiesId[i], table)
			if dbErr.Err != nil {
				rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
				return
			}
		}

		if filterParams {
			index := 1
			var membershipCheck bool
			for index < len(ps) {
				if dbTables[index] != database.LikeTable {
					membershipCheck, dbErr = rt.db.DoesEntityBelongsTo(entitiesId[index], entitiesId[index-1], dbTables[index])
					if !membershipCheck {
						if dbErr.Err != nil {
							rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
						} else {
							rt.LoggerAndHttpErrorSender(w, errors.New("entity does not belong to the parent"), utils.HttpError{StatusCode: 400})
						}
						return
					}
					index++
				}
			}
		}

		authorizationHeader := r.Header.Get("Authorization")
		token := utils.GetAuthenticationToken(authorizationHeader)
		if !token.IsValid() {
			rt.LoggerAndHttpErrorSender(w, errors.New("token invalid"), utils.HttpError{StatusCode: 401})
			return
		}

		fn(w, r, ps, token)

		// Call the next handler in chain (usually, the handler function for the path)
	}
}

func (rt *_router) authWrap(fn httpRouterHandler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {

		isAuthorized := utils.Authorize(token, ps.ByName("user_id"))

		if !isAuthorized {
			rt.LoggerAndHttpErrorSender(w, errors.New("user not authorized"), utils.HttpError{StatusCode: 401})
			return
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, token)
	}
}
