package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto/service/database"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var user *database.User
	user, err := rt.db.GetUser("Lorenzo")

	if err != nil {
		rt.baseLogger.Errorf("error getting user: %v", err)
		return
	} else {
		err = json.NewEncoder(w).Encode(*user)
		if err != nil {
			rt.baseLogger.Errorf("error encoding user: %v", err)
			return
		}
	}
}
