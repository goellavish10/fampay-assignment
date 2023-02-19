package main

import (
	"log"
	"net/http"

	"github.com/goellavish10/fampay-assignment/configs"
	"github.com/goellavish10/fampay-assignment/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Connect to database
	configs.ConnectDB()

	// Routes
	routes.SearchRoutes(router)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
