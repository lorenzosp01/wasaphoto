package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type JSONErrorMessage struct {
	Message string `json:"message"`
}

func (rt *_router) createFountain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "application/json")

	var fountain Fountain
	err := json.NewDecoder(r.Body).Decode(&fountain)

	fmt.Println(fountain)
	if err != nil {
		rt.baseLogger.WithError(err).Warningf("createFountain: error decoding request body")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMessage{Message: "error decoding request body"})
		fmt.Println(err)
	}

}
