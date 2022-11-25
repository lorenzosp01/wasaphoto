package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Fountain struct {
	Id        int    `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Status    string `json:"status"`
}

func (rt *_router) getFountains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	f1 := Fountain{1, "10", "20", "good"}
	f2 := Fountain{2, "11", "21", "faulty"}
	f3 := Fountain{3, "11", "21", "faulty"}
	fountains := []Fountain{f1, f2, f3}

	// OPPURE json.NewEncoder(w).Encode(string(jsonString))
	jsonString, err := json.Marshal(fountains)

	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jsonString)
	if err != nil {
		fmt.Println(err)
	}
}
