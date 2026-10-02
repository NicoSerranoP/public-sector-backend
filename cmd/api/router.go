package main

import (
	"net/http"
	"public-sector-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func limitBodySize(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1048576) // 1MB limit
	c.Next()
}

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(limitBodySize)

	r.GET("/health", handlers.HealthHandler)

	r.POST("/voter", handlers.RegisterVoterHandler)

	r.GET("/voter/:id", handlers.GetVoterHandler)

	r.POST("/blacklist", handlers.RegisterBlacklistVoterHandler)

	r.GET("/blacklist", handlers.GetBlacklistedVotersHandler)

	r.GET("/blacklist/:address", handlers.GetIsVoterBlacklistedHandler)

	r.DELETE("/blacklist/:address", handlers.RemoveVoterFromBlacklistHandler)

	return r
}
