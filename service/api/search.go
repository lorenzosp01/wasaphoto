package api

import (
	"encoding/json"
	"net/http"
	"wasaphoto/service/database"
)

func (rt *_router) doSearch(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	pattern := r.URL.Query().Get("pattern")
	authUserId := params["token"]

	dbUsers, dbErr := rt.db.DoSearch(pattern)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	var users []User
	for _, dbUser := range dbUsers {
		var user User
		user.fromDatabase(dbUser)
		isBanned, dbErr := rt.db.IsUserTargeted(user.Id, authUserId, database.BanTable)
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		}
		if !isBanned {
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		users = make([]User, 0)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userList{users})
}
