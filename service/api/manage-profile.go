package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	userId := params["token"]

	photo, err := io.ReadAll(r.Body)
	if err != nil || len(photo) == 0 {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.InsertPhoto(photo, userId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo uploaded successfully"))
}

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, params map[string]int64) {

	photoId := params["photo_id"]
	authUserId := params["token"]
	userId := params["user_id"]

	userIsBanned, dbErr := rt.db.IsUserTargeted(userId, authUserId, database.BanTable)
	if userIsBanned {
		dbErr.CustomMessage = utils.BannedMessage
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	} else {
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		}
	}

	var image []byte
	image, dbErr = rt.db.GetImage(photoId, userId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	userId := params["user_id"]

	var newUsername Username
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400, Message: "Invalid request body"})
		return
	}

	if !newUsername.IsValid() {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400, Message: "Invalid username"})
		return
	}

	dbErr := rt.db.ChangeUsername(userId, newUsername.Username)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	user := User{
		Id:       userId,
		Username: newUsername.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	photoId := params["photo_id"]
	userId := params["user_id"]

	dbErr := rt.db.DeletePhoto(photoId, userId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo deleted successfully"))
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, params map[string]int64) {

	authUserId := params["token"]
	userId := params["user_id"]

	var userProfile UserProfile
	// get query params
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400, Message: "Query paramaters badly formatted"})
		return
	}

	amount, err := strconv.ParseInt(r.URL.Query().Get(r.URL.Query().Get("amount")), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400, Message: "Query paramaters badly formatted"})
		return
	}

	var userIsBanned bool
	var dbErr database.DbError
	userIsBanned, dbErr = rt.db.IsUserTargeted(userId, authUserId, database.BanTable)
	if userIsBanned {
		dbErr.CustomMessage = utils.BannedMessage
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	} else {
		if dbErr.InternalError != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
			return
		}
	}

	up, dbErr := rt.db.GetUserProfile(userId, amount, offset)
	userProfile.fromDatabase(up)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userProfile)

}
