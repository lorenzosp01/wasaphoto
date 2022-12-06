package api

import (
	"encoding/json"
	"net/http"
	"wasaphoto/service/database"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId, _ := params["token"]

	dbUsers, dbErr := rt.db.GetUsersList(authUserId, database.FollowTable)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	var followedUsers []User
	for _, dbUser := range dbUsers {
		var user User
		user.fromDatabase(dbUser)
		followedUsers = append(followedUsers, user)
	}

	var photos []Photo
	for _, user := range followedUsers {
		dbPhotos, dbErr := rt.db.GetUserPhotos(user.Id, params["amount"], params["offset"])
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}

		for _, dbPhoto := range dbPhotos {
			var photo Photo
			photo.fromDatabase(dbPhoto)
			photos = append(photos, photo)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserStream{photos})

}
