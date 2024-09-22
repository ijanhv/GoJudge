package main

import (
	"gojudge/db"
	"gojudge/routes"
)

func main() {
	db.InitDB()

	routes.StartServer()
}
