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
		isTargeted, dbErr := rt.db.IsUserTargeted(user.Id, authUserId, database.BanTable)
		if !isTargeted {
			if dbErr.InternalError != nil {
				rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
				return
			}
			users = append(users, user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userList{users})
}
