package main

import (
	"log"

	"public-sector-backend/internal/config"
)

func main() {
	r := NewRouter()

	var address string

	if config.IS_IN_PRODUCTION {
		address = "0.0.0.0:8080"
	} else {
		address = "127.0.0.1:8080"
	}

	if err := r.Run(address); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
