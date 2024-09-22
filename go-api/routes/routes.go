package routes

import (
	"gojudge/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/submission", controllers.Submission).Methods("POST")
	log.Println("Server running in port 8001 ")
    log.Fatal(http.ListenAndServe(":8001", router))
}
