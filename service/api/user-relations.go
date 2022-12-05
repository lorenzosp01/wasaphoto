package api

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

type userList struct {
	Users []User `json:"users"`
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	rt.targetUser(w, r, ps, token, database.BanTable)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	rt.targetUser(w, r, ps, token, database.FollowTable)
}

func (rt *_router) targetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token, entityTable string) {
	authUserId, err := strconv.ParseInt(token.Value, 10, 64)
	targetedUserId, err := strconv.ParseInt(ps.ByName("targeted_user_id"), 10, 64)

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 403})
		return
	}

	dbErr := rt.db.TargetUser(authUserId, targetedUserId, entityTable)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	switch entityTable {
	case database.BanTable:
		_, _ = w.Write([]byte("User banned successfully"))
	case database.FollowTable:
		_, _ = w.Write([]byte("User followed successfully"))
	}

}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	rt.untargetUser(w, r, ps, token, database.BanTable)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	rt.untargetUser(w, r, ps, token, database.FollowTable)
}

func (rt *_router) untargetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token, entityTable string) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	targetedUserId, _ := strconv.ParseInt(ps.ByName("targeted_user_id"), 10, 64)

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("a user can't untarget himself"), utils.HttpError{StatusCode: 403})
		return
	}

	dbErr := rt.db.UntargetUser(authUserId, targetedUserId, entityTable)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	switch entityTable {
	case database.BanTable:
		_, _ = w.Write([]byte("User unbanned successfully"))
	case database.FollowTable:
		_, _ = w.Write([]byte("User unfollowed successfully"))
	}

}

func (rt *_router) getFollowedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	rt.getUsersList(w, r, ps, token, database.FollowTable)
}

func (rt *_router) getBannedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	rt.getUsersList(w, r, ps, token, database.BanTable)
}

func (rt *_router) getUsersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token, entityTable string) {
	userId, _ := strconv.ParseInt(token.Value, 10, 64)

	dbUsers, dbErr := rt.db.GetUsersList(userId, entityTable)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	var users []User
	for _, dbUser := range dbUsers {
		var user User
		user.fromDatabase(dbUser)
		users = append(users, user)
	}

	userList := userList{
		Users: users,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userList)
}
