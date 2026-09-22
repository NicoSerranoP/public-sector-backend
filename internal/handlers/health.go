package handlers

import (
	"net/http"
	"public-sector-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "healthy",
		"APP_ENV": config.APP_ENV,
	})
}
