package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadConfig: loads env variables
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// Get: returns value for key with fallback
func Get(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
