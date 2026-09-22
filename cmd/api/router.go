package main

import (
	"public-sector-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.HealthHandler)

	return r
}
