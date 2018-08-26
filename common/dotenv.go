package common

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv values from .env file into environment variables
func LoadEnv() {
	err := godotenv.Load("/etc/logstv/.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

// GetEnv retrieve a value from environment variables
func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Missing env var: " + key)
}
