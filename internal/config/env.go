package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var APP_ENV string
var IS_IN_PRODUCTION bool

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	APP_ENV = os.Getenv("APP_ENV")
	IS_IN_PRODUCTION = APP_ENV == "production"
}
