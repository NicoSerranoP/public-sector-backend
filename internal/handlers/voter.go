package handlers

import (
	"log"
	"net/http"
	"time"

	"public-sector-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func RegisterVoterHandler(c *gin.Context) {
	var body RegisterVoterDTO

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: validate the user is sending his own identity number and not someone's else

	_, err := database.Get().Exec(
		c.Request.Context(),
		"INSERT INTO voters (id, address) VALUES ($1, $2)",
		body.Id, body.Address,
	)

	if err != nil {
		message := "failed to register voter"
		log.Println(message+":", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      body.Id,
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

	var voter struct {
		Id        string    `json:"id"`
		Address   string    `json:"address"`
		CreatedAt time.Time `json:"created_at"`
	}

	if err := row.Scan(&voter.Id, &voter.Address, &voter.CreatedAt); err != nil {
		log.Println("failed to get voter:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "voter not found"})
		return
	}

	c.JSON(http.StatusOK, voter)
}

func GetBlacklistedVotersHandler(c *gin.Context) {
	// TODO: implement this
}
