package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/utils"
)

// todo capire se serve controllare l'esistenza dell'utente nel db
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	authUserId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))
	if httpErr.StatusCode != 0 {
		rt.LoggerAndHttpErrorSender(w, err, httpErr)
		return
	}

	targetedUserId, err := strconv.ParseInt(ps.ByName("targeted_user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 403})
		return
	}
	//todo capire se quel if si può fare nella funzione
	userIsBanned, dbErr := rt.db.IsUserBannedBy(targetedUserId, authUserId)
	if userIsBanned {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		} else {
			rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 409})
		}
		return
	}

	dbErr = rt.db.BanUser(authUserId, targetedUserId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("User banned successfully"))
}
