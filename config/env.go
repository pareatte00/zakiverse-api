package config

import (
	"log"

	"github.com/joho/godotenv"
)

func loadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load environment file: %v", err)
	}
}
