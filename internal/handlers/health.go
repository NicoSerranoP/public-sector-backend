package handlers

import (
	"net/http"
	"public-sector-backend/internal/config"
	"public-sector-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	healthy := true
	message := "healthy"

	successfulQuery := database.Get().Ping(c.Request.Context())

	if successfulQuery != nil {
		healthy = false
		message = "could not ping the database"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"healthy": healthy,
		"APP_ENV": config.APP_ENV,
	})
}
