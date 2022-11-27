package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var username utils.Username
	err := json.NewDecoder(r.Body).Decode(&username)

	if err != nil {
		// Request body is not a valid JSON
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid JSON"))
		rt.baseLogger.Errorf("error decoding username: %v", err)
		return
	} else {
		if !username.IsValid() {
			// Username is not valid
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Invalid username"))
			return
		}
	}

	var id utils.UserIdentifier
	var dbErr database.DbError

	id.Id, dbErr = rt.db.GetUserId(username.Username)
	// if an error occurred while getting the user id
	if dbErr.Err != nil {
		httpErr := dbErr.ToHttp()
		w.WriteHeader(httpErr.StatusCode)
		_, _ = w.Write([]byte(httpErr.Message))
		rt.baseLogger.WithError(err).Error("error getting user id")
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(id)

}
