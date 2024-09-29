package main

import (
	"gojudge/db"
	"gojudge/routes"
	"gojudge/storage"
	"log"
	"github.com/joho/godotenv"
)


func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}

func main() {
	
	 db.InitDB()
	storage.InitStorage()

	routes.StartServer()
}
