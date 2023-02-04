package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

type userList struct {
	Users []User `json:"users"`
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	rt.targetUser(w, params, database.BanTable)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	rt.targetUser(w, params, database.FollowTable)
}

func (rt *_router) targetUser(w http.ResponseWriter, params map[string]int64, entityTable string) {

	authUserId := params["token"]
	targetedUserId := params["targeted_user_id"]

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("a user can't target himself"), utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't target yourself"})
		return
	}

	if entityTable == database.FollowTable {
		// Check if the targeted user id banned the authenticated one
		isBanned, dbErr := rt.db.IsUserTargeted(targetedUserId, authUserId, database.BanTable)
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		}

		if isBanned {
			rt.LoggerAndHttpErrorSender(w, nil, utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't follow who banned you"})
			return
		}
	} else if entityTable != database.BanTable {
		return
	}

	isOperationSuccessful, dbErr := rt.db.TargetUser(authUserId, targetedUserId, entityTable)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	if !isOperationSuccessful {
		rt.LoggerAndHttpErrorSender(w, nil, utils.HttpError{StatusCode: http.StatusInternalServerError, Message: "Can't target user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	switch entityTable {
	case database.BanTable:
		_, _ = w.Write([]byte("User banned successfully"))
	case database.FollowTable:
		_, _ = w.Write([]byte("User followed successfully"))
	}

}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	rt.untargetUser(w, params, database.BanTable)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	rt.untargetUser(w, params, database.FollowTable)
}

func (rt *_router) untargetUser(w http.ResponseWriter, params map[string]int64, entityTable string) {
	authUserId := params["token"]
	targetedUserId := params["targeted_user_id"]

	if targetedUserId == authUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("a user can't untarget himself"), utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't untarget yourself"})
		return
	}

	if entityTable == database.FollowTable {
		isTargeted, dbErr := rt.db.IsUserTargeted(targetedUserId, authUserId, database.BanTable)
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		}
		if isTargeted {
			rt.LoggerAndHttpErrorSender(w, nil, utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't unfollow who banned you"})
			return
		}
	} else if entityTable != database.BanTable {
		return
	}

	isOperationSuccessful, dbErr := rt.db.UntargetUser(authUserId, targetedUserId, entityTable)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	if !isOperationSuccessful {
		rt.LoggerAndHttpErrorSender(w, nil, utils.HttpError{StatusCode: http.StatusNotFound, Message: "Target requested not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	switch entityTable {
	case database.BanTable:
		_, _ = w.Write([]byte("User unbanned successfully"))
	case database.FollowTable:
		_, _ = w.Write([]byte("User unfollowed successfully"))
	}

}

func (rt *_router) getFollowedUsers(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	rt.getUsersList(w, r, params, database.FollowTable)
}

func (rt *_router) getBannedUsers(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	rt.getUsersList(w, r, params, database.BanTable)
}

func (rt *_router) getUsersList(w http.ResponseWriter, r *http.Request, params map[string]int64, entityTable string) {
	userId := params["user_id"]

	dbUsers, dbErr := rt.db.GetUsersList(userId, entityTable)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	var users []User
	if dbUsers == nil {
		users = make([]User, 0)
		return
	}

	for _, dbUser := range dbUsers {
		var user User
		user.fromDatabase(dbUser)
		users = append(users, user)
	}

	userList := userList{
		Users: users,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userList)
}
