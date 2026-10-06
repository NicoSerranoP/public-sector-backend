package handlers

import (
	"log"
	"net/http"

	"public-sector-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func RegisterVoterHandler(c *gin.Context) {
	var body RegisterVoterDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: validate the user is sending his own identity number and not someone's else

	_, err := database.Get().Exec(
		c.Request.Context(),
		"INSERT INTO voters (id, address) VALUES ($1, $2)",
		body.ID, body.Address,
	)

	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register voter"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      body.ID,
		"address": body.Address,
	})
}

func GetVoterHandler(c *gin.Context) {
	id := c.Param("id")
	idType := c.DefaultQuery("type", "id")

	if idType != "id" && idType != "address" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong type of id"})
		return
	}

	row := database.Get().QueryRow(
		c.Request.Context(),
		"SELECT * FROM voters WHERE "+idType+"=$1",
		id,
	)

	var voter database.VoterType

	if err := row.Scan(&voter.ID, &voter.Address, &voter.CreatedAt); err != nil {
		log.Println("failed to get voter:", err)

		c.JSON(http.StatusNotFound, gin.H{"error": "voter not found"})
		return
	}

	c.JSON(http.StatusOK, voter)
}
