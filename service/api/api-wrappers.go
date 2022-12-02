package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/utils"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
//func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
//	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//		reqUUID, err := uuid.NewV4()
//		if err != nil {
//			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		var ctx = reqcontext.RequestContext{
//			ReqUUID: reqUUID,
//		}
//
//		// Create a request-specific logger
//		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
//			"reqid":     ctx.ReqUUID.String(),
//			"remote-ip": r.RemoteAddr,
//		})
//
//		// Call the next handler in chain (usually, the handler function for the path)
//		fn(w, r, ps, ctx)
//	}
//}

// Parses the request and checks if path in the id are valid and entity with that exists
func (rt *_router) wrap(fn httpRouterHandler, dbTables []string) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		for i, table := range dbTables {
			entityId, err := strconv.ParseInt(ps[i].Value, 10, 64)
			if err != nil {
				rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
				return
			}

			dbErr := rt.db.EntityExists(entityId, table)
			if dbErr.Err != nil {
				rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
				return
			}
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps)
	}
}

func (rt *_router) authWrap(fn httpRouterHandler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		authorizationHeader := r.Header.Get("Authorization")
		token := utils.GetAuthenticationToken(authorizationHeader)
		if !token.IsValid() {
			rt.LoggerAndHttpErrorSender(w, errors.New("token invalid"), utils.HttpError{StatusCode: 401})
			return
		}

		isAuthorized := utils.Authorize(token, ps[0].Value)

		if !isAuthorized {
			rt.LoggerAndHttpErrorSender(w, errors.New("user not authorized"), utils.HttpError{StatusCode: 401})
			return
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps)
	}
}
