package handlers

import (
	"log"
	"net/http"

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
	// TODO: implement this
	/*
		row := database.Get().QueryRow(c.Request.Context(), "SELECT * FROM voters WHERE id=$1", c.Param("id"))

		return row
	*/
}

func GetBlacklistedVotersHandler(c *gin.Context) {
	// TODO: implement this
}
