package main

import (
	"gojudge/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	routes.RegisterRoutes(router)

	log.Println("Server running on port 8001..")
	log.Fatal(http.ListenAndServe(":8000", router))

}
