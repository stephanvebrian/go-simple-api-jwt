package home

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func New(router *mux.Router) *mux.Router {
	router.HandleFunc("/", indexHandler).Methods("GET")

	return router
}

type response struct {
	Message string `json:"message"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	res := response{
		Message: "Application Up",
	}

	json.NewEncoder(w).Encode(res)
}
