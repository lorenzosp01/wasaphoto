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
	ownerId, err := strconv.Atoi(ps.ByName("user_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid user id"))
		rt.baseLogger.WithError(err).Error("error parsing user id")
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

	err = rt.db.InsertPhoto(photo, int64(ownerId))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error inserting photo"))
		rt.baseLogger.WithError(err).Error("error inserting photo in the database")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo uploaded successfully"))
}

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	photo_id, err := strconv.Atoi(ps.ByName("photo_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid photo id"))
		rt.baseLogger.WithError(err).Error("error parsing photo id")
		return
	}

	image, err := rt.db.GetImage(int64(photo_id))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting image"))
		rt.baseLogger.WithError(err).Error("error getting image from the database")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}
