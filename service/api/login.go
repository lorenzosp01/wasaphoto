package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto/service/database"
)

// todo change response and logger
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var username Username
	err := json.NewDecoder(r.Body).Decode(&username)

	if err != nil {
		// Request body is not a valid JSON
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid JSON"))
		rt.baseLogger.Errorf("error decoding username: %v", err)
		return
	}

	if !username.IsValid() {
		// Username is not valid
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid username"))
		return
	}

	var id UserIdentifier
	var dbErr database.DbError
	var created bool = false

	// modificare il doc di OpenAPI
	created, id.Id, dbErr = rt.db.GetUserId(username.Username)
	// if an error occurred while getting the user id
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	//todo settare gli header ovunque
	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-type", "application/json")
	_ = json.NewEncoder(w).Encode(id)
}
