package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadSecret() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

}
