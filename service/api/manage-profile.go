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
	userId, _ := params["token"]

	photo, err := io.ReadAll(r.Body)
	if err != nil || len(photo) == 0 {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.InsertPhoto(photo, userId)
	if dbErr.Err != nil {
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

	userIsBanned, dbErr := rt.db.IsUserAlreadyTargeted(authUserId, userId, database.BanTable)
	if userIsBanned {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}

	var image []byte
	image, dbErr = rt.db.GetImage(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
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
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	if !newUsername.IsValid() {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.ChangeUsername(userId, newUsername.Username)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	user := User{
		Id:       userId,
		Username: newUsername.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	photoId := params["photo_id"]

	var dbErr database.DbError

	dbErr = rt.db.DeletePhoto(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo deleted successfully"))
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, params map[string]int64) {

	authUserId := params["token"]
	userId := params["user_id"]

	photosAmount, err := strconv.ParseInt(r.URL.Query().Get("amount"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	photosOffset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	var userIsBanned bool
	var dbErr database.DbError
	userIsBanned, dbErr = rt.db.IsUserAlreadyTargeted(authUserId, userId, database.BanTable)
	if userIsBanned {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}

	var userProfile UserProfile
	up, dbErr := rt.db.GetUserProfile(userId, photosAmount, photosOffset)
	userProfile.fromDatabase(up)

	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userProfile)

}
