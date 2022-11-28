package api

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"wasaphoto/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))

	if httpErr.Message != "" {
		w.WriteHeader(httpErr.StatusCode)
		_, _ = w.Write([]byte(httpErr.Message))
		rt.baseLogger.WithError(errors.New(httpErr.Message)).Error("error authorizing user")
		return
	}

	var photo []byte
	photo, err = io.ReadAll(r.Body)
	if err != nil || len(photo) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid photo"))
		rt.baseLogger.WithError(err).Error("error reading photo")
		return
	}

	dbErr = rt.db.InsertPhoto(photo, userId)
	if dbErr.Err != nil {
		httpErr = dbErr.ToHttp()
		w.WriteHeader(httpErr.StatusCode)
		_, _ = w.Write([]byte(httpErr.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error inserting photo in the database")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo uploaded successfully"))
}

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	photoId, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid photo id"))
		rt.baseLogger.WithError(err).Error("error parsing photo id")
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get user")
		return
	}

	var image []byte
	image, dbErr = rt.db.GetImage(photoId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get image")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get user")
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))

	if httpErr.Message != "" {
		w.WriteHeader(httpErr.StatusCode)
		_, _ = w.Write([]byte(httpErr.Message))
		rt.baseLogger.WithError(errors.New(httpErr.Message)).Error("error authorizing user")
		return
	}

	var newUsername utils.Username
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid JSON"))
		rt.baseLogger.WithError(err).Error("error decoding username")
		return
	}

	if !newUsername.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid username"))
		rt.baseLogger.WithError(err).Error("error validating username")
		return
	}

	dbErr = rt.db.ChangeUsername(userId, newUsername.Username)
	if dbErr.Err != nil {
		httpErr = dbErr.ToHttp()
		w.WriteHeader(httpErr.StatusCode)
		_, _ = w.Write([]byte(httpErr.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error changing username")
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
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	photoId, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid photo id"))
		rt.baseLogger.WithError(err).Error("error parsing photo id")
		return
	}

	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get user")
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))

	if httpErr.Message != "" {
		w.WriteHeader(httpErr.StatusCode)
		_, _ = w.Write([]byte(httpErr.Message))
		rt.baseLogger.WithError(errors.New(httpErr.Message)).Error("error authorizing user")
		return
	}

	var isPhotoOwner bool
	// Check if photo exists and if it belongs to the user
	isPhotoOwner, dbErr = rt.db.IsPhotoOwner(photoId, userId)
	if !isPhotoOwner {
		httpErr = dbErr.ToHttp()
		if httpErr.StatusCode != http.StatusInternalServerError {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("You cannot delete the photo"))
			rt.baseLogger.WithError(dbErr.Err).Error("The user cannot access to that photo")
		} else {
			w.WriteHeader(httpErr.StatusCode)
			_, _ = w.Write([]byte(httpErr.Message))
			rt.baseLogger.WithError(dbErr.Err).Error("Server error")
		}
		return
	}

	dbErr = rt.db.DeletePhoto(photoId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get image")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo deleted successfully"))
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId, err := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	//todo capire se è necessario
	dbErr := rt.db.UserExists(userId)
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get user")
		return
	}

	token := utils.GetAuthenticationToken(r.Header.Get("Authorization"))
	if !token.IsValid() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("You are not authorized to access this resource"))
		rt.baseLogger.WithError(errors.New("error authorizing user")).Error("error authorizing user")
		return
	}
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)

	var userIsBanned bool
	userIsBanned, dbErr = rt.db.IsUserBannedBy(authUserId, userId)
	if userIsBanned {
		if dbErr.Err != nil {
			httpErr := dbErr.ToHttp()
			w.WriteHeader(httpErr.StatusCode)
			rt.baseLogger.WithError(dbErr.Err).Error(httpErr.Message)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Banned by that user"))
		rt.baseLogger.Errorf("The user is banned")
		return
	}

	var userProfile UserProfile
	up, dbErr := rt.db.GetUserProfile(userId)
	userProfile.fromDatabase(up)

	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get user profile")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userProfile)

}
