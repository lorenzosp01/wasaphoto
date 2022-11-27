package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"wasaphoto/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, err := strconv.Atoi(ps.ByName("user_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	dbErr := rt.db.UserExists(int64(userId))
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	httpErr := utils.Authorize(authorizationHeader, ps.ByName("user_id"))

	if httpErr != (utils.HttpError{}) {
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

	dbErr = rt.db.InsertPhoto(photo, int64(userId))
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error inserting photo in the database")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo uploaded successfully"))
}

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
		return
	}

	photoId, err := strconv.Atoi(ps.ByName("photo_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid photo id"))
		rt.baseLogger.WithError(err).Error("error parsing photo id")
		return
	}

	dbErr := rt.db.UserExists(int64(userId))
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get user")
		return
	}

	var image []byte
	image, dbErr = rt.db.GetImage(int64(photoId))
	if dbErr.Err != nil {
		httpError := dbErr.ToHttp()
		w.WriteHeader(httpError.StatusCode)
		_, _ = w.Write([]byte(httpError.Message))
		rt.baseLogger.WithError(dbErr.Err).Error("error during get image")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}
