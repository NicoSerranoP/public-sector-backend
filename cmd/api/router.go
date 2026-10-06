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
	router := gin.Default()
	router.SetTrustedProxies(nil) // set nginx proxy IP if using reverse proxy

	router.Use(limitBodySize)

	router.GET("/health", handlers.HealthHandler)

	router.POST("/voter", handlers.RegisterVoterHandler)
	router.GET("/voter/:id", handlers.GetVoterHandler)

	router.POST("/blacklist", handlers.RegisterBlacklistVoterHandler)
	router.GET("/blacklist", handlers.GetBlacklistedVotersHandler)
	router.GET("/blacklist/:address", handlers.GetIsVoterBlacklistedHandler)
	router.DELETE("/blacklist/:address", handlers.RemoveVoterFromBlacklistHandler)

	router.POST("/relay/vote", handlers.RelayCastVoteHandler)
	router.POST("/relay/proposal", handlers.RelayCreateProposalHandler)

	return router
}
