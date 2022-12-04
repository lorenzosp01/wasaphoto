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

// Returns a photo per user followed by the authenticated user
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)

	photosAmount, err := strconv.ParseInt(r.URL.Query().Get("amount"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, errors.New("QueryString badly formatted"), utils.HttpError{StatusCode: 400})
		return
	}

	photosOffset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, errors.New("QueryString badly formatted"), utils.HttpError{StatusCode: 400})
		return
	}

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
		dbPhotos, dbErr := rt.db.GetUserPhotos(user.Id, photosAmount, photosOffset)
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
