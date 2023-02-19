package routes

import "github.com/gorilla/mux"

func SearchRoutes(router *mux.Router) {
	router.HandleFunc("/search", nil).Methods("GET")
}
