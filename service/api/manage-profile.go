package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	userId, _ := strconv.ParseInt(token.Value, 10, 64)

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

//todo controllare h
func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {

	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)

	var dbErr database.DbError
	var image []byte
	image, dbErr = rt.db.GetImage(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, p httprouter.Params, token utils.Token) {
	userId, _ := strconv.ParseInt(token.Value, 10, 64)

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

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	userId, _ := strconv.ParseInt(token.Value, 10, 64)
	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)

	var dbErr database.DbError
	var isPhotoOwner bool
	// Check if photo exists and if it belongs to the user
	isPhotoOwner, dbErr = rt.db.IsPhotoOwner(photoId, userId)
	if !isPhotoOwner {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	dbErr = rt.db.DeletePhoto(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo deleted successfully"))
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {

	userId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)

	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)

	//todo spostare if dentro la funzione di interazione col db
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
	up, dbErr := rt.db.GetUserProfile(userId)
	userProfile.fromDatabase(up)

	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userProfile)

}
