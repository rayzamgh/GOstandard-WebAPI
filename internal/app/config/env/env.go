package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Get(name string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	val := os.Getenv(name)
	return val
}
