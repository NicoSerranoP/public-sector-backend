package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//revive:disable:var-naming
var APP_ENV string
var DATABASE_URL string
var IS_IN_PRODUCTION bool

//revive:enable:var-naming

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	APP_ENV = os.Getenv("APP_ENV")
	DATABASE_URL = os.Getenv("DATABASE_URL")
	IS_IN_PRODUCTION = APP_ENV == "production"
}
