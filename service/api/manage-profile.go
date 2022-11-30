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

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))
	if httpErr.StatusCode != 0 {
		rt.LoggerAndHttpErrorSender(w, err, httpErr)
		return
	}

	var photo []byte
	photo, err = io.ReadAll(r.Body)
	if err != nil || len(photo) == 0 {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr = rt.db.InsertPhoto(photo, userId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo uploaded successfully"))
}

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	photoId, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	var image []byte
	image, dbErr = rt.db.GetImage(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))

	if httpErr.StatusCode != 0 {
		rt.LoggerAndHttpErrorSender(w, err, httpErr)
		return
	}

	var newUsername utils.Username
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	if !newUsername.IsValid() {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr = rt.db.ChangeUsername(userId, newUsername.Username)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	user := User{
		Id:       userId,
		Username: newUsername.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	photoId, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))

	if httpErr.Message != "" {
		rt.LoggerAndHttpErrorSender(w, err, httpErr)
		return
	}

	var isPhotoOwner bool
	// Check if photo exists and if it belongs to the user
	isPhotoOwner, dbErr = rt.db.IsPhotoOwner(photoId, userId)
	if !isPhotoOwner {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	dbErr = rt.db.DeletePhoto(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo deleted successfully"))
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	//todo capire se è necessario
	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	token := utils.GetAuthenticationToken(r.Header.Get("Authorization"))
	if !token.IsValid() {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	var userIsBanned bool
	userIsBanned, dbErr = rt.db.IsUserAlreadyTargeted(authUserId, userId, database.BanTable)
	if userIsBanned {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
			return
		}
	}

	var userProfile UserProfile
	up, dbErr := rt.db.GetUserProfile(userId)
	userProfile.fromDatabase(up)

	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userProfile)

}
