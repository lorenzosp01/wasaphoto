package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto/service/utils"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var username utils.Username
	err := json.NewDecoder(r.Body).Decode(&username)

	if err != nil {
		rt.baseLogger.Errorf("error decoding username: %v", err)
		return
	} else {
		var userId utils.UserIdentifier
		userId.Id, err = rt.db.GetUserId(username.Username)
		if err != nil {
			rt.baseLogger.Errorf("error getting user id: %v", err)
			return
		} else {
			err = json.NewEncoder(w).Encode(userId)
			if err != nil {
				rt.baseLogger.Errorf("error encoding user: %v", err)
				return
			}
		}
	}
}
