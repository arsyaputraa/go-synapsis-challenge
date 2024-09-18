package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value from key ---
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

func LoadEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}
