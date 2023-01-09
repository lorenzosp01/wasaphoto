package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"wasaphoto/service/utils"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]

	var photos []Photo

	// todo implementare meglio la restituzione dello stream
	dbPhotos, dbErr := rt.db.GetMyStream(authUserId, params["offset"], params["amount"])
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	if len(dbPhotos) == 0 {
		httpErr := utils.HttpError{StatusCode: 404, Message: "No photos found"}
		rt.LoggerAndHttpErrorSender(w, errors.New("no photos for that user stream"), httpErr)
		return
	}

	for _, dbPhoto := range dbPhotos {
		var photo Photo
		photo.fromDatabase(dbPhoto)
		photos = append(photos, photo)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserStream{photos})

}
