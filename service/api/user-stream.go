package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wasaphoto/service/utils"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	var photos []Photo

	// get query params
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: http.StatusBadRequest, Message: "Query paramaters badly formatted"})
		return
	}

	amount, err := strconv.ParseInt(r.URL.Query().Get(r.URL.Query().Get("amount")), 10, 64)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: http.StatusBadRequest, Message: "Query paramaters badly formatted"})
		return
	}

	dbPhotos, dbErr := rt.db.GetMyStream(authUserId, offset, amount)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	if len(dbPhotos) == 0 {
		httpErr := utils.HttpError{StatusCode: http.StatusNotFound, Message: "No photos found"}
		rt.LoggerAndHttpErrorSender(w, errors.New("no photos for that user stream"), httpErr)
		return
	}

	for _, dbPhoto := range dbPhotos {
		var photo Photo
		photo.fromDatabase(dbPhoto)
		photos = append(photos, photo)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserStream{photos})

}
