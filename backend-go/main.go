package main

import (
	"backend/database"
	"backend/handler"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Initiating...")
	db := database.InitDB("./database/shifts.db")

	//internal module
	http.HandleFunc("/healthz", handler.HandleHealth)
	http.Handle("/login", handler.CORS(handler.LoginHandler(db)))

	//list of handler
	handler.RegisterEmployeeRoutes(db)
	handler.RegisterAdminRoutes(db)

	log.Println("Server running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
