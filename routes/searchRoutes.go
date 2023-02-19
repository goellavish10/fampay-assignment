package routes

import (
	"github.com/goellavish10/fampay-assignment/controllers"
	"github.com/gorilla/mux"
)

func SearchRoutes(router *mux.Router) {
	router.HandleFunc("/search", controllers.SearchController()).Methods("GET")
}
