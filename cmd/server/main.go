package main

import (
	"log"
	"net/http"

	"github.com/Employes-Side/employee-side/internal/config"
	"github.com/Employes-Side/employee-side/internal/db"
	"github.com/Employes-Side/employee-side/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	cfg := config.LoadConfig()

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect with database: %v", err)
	}

	defer dbConn.Close()

	http.HandleFunc("/users", handlers.UserHandler(dbConn))

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
