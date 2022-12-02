package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rt.targetUser(w, r, ps, database.BanTable)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rt.targetUser(w, r, ps, database.FollowTable)
}

// todo per il follow stare attento a controllare che l'utente non possa seguire un utente che lo ha bannato e vicersa
func (rt *_router) targetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, tableToUse string) {
	authUserId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	targetedUserId, err := strconv.ParseInt(ps.ByName("targeted_user_id"), 10, 64)

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 403})
		return
	}

	//todo i possibili doppi inserimenti di riga nel database vanno gestiti come dberror
	isTargeted, dbErr := rt.db.IsUserAlreadyTargeted(authUserId, targetedUserId, tableToUse)
	if isTargeted {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		} else {
			rt.LoggerAndHttpErrorSender(w, errors.New("user already targeted"), utils.HttpError{StatusCode: 409})
		}
		return
	}

	dbErr = rt.db.TargetUser(authUserId, targetedUserId, tableToUse)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	switch tableToUse {
	case database.BanTable:
		_, _ = w.Write([]byte("User banned successfully"))
	case database.FollowTable:
		_, _ = w.Write([]byte("User followed successfully"))
	}

}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rt.untargetUser(w, r, ps, database.BanTable)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rt.untargetUser(w, r, ps, database.FollowTable)
}

func (rt *_router) untargetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, tableToUse string) {
	authUserId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	targetedUserId, _ := strconv.ParseInt(ps.ByName("targeted_user_id"), 10, 64)

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("a user can't untarget himself"), utils.HttpError{StatusCode: 403})
		return
	}

	dbErr := rt.db.UntargetUser(authUserId, targetedUserId, tableToUse)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	switch tableToUse {
	case database.BanTable:
		_, _ = w.Write([]byte("User unbanned successfully"))
	case database.FollowTable:
		_, _ = w.Write([]byte("User unfollowed successfully"))
	}

}
