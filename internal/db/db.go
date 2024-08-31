package db

import (
	"database/sql"
	"log"

	"github.com/Employes-Side/employee-side/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the Mysql database connection")
	return db, nil
}
