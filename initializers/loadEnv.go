package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Could not load .env file")
	}
}
