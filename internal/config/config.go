package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABSE_URL is required")
	}

	return &Config{
		DatabaseURL: dbURL,
	}
}
