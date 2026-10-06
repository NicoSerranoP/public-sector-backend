package main

import (
	"context"
	"log"

	"public-sector-backend/internal/config"
	"public-sector-backend/internal/database"
)

func main() {
	if err := database.Initialize(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer database.Get().Close()

	r := NewRouter()

	var address string

	if config.IS_IN_PRODUCTION {
		address = "0.0.0.0:8080"
	} else {
		address = "127.0.0.1:8080"
	}

	if err := r.Run(address); err != nil {
		log.Printf("failed to run server: %v", err)
		return
	}
}
