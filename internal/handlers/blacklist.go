package handlers

import (
	"log"
	"net/http"
	"public-sector-backend/internal/database"

	"github.com/gin-gonic/gin"
)

func RegisterBlacklistVoterHandler(c *gin.Context) {
	var body BlacklistVoterDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.Get().Exec(
		c.Request.Context(),
		"INSERT INTO blacklist (address, duration) VALUES ($1, make_interval(mins => $2))",
		body.Address,
		body.Duration,
	)

	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to blacklist voter"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"address":  body.Address,
		"duration": body.Duration,
	})
}

func GetBlacklistedVotersHandler(c *gin.Context) {
	rows, err := database.Get().Query(
		c.Request.Context(),
		`SELECT address, (EXTRACT(EPOCH FROM duration) / 60)::int, is_active, created_at
		FROM blacklist WHERE is_active = TRUE`,
	)

	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get blacklisted voters"})
		return
	}

	defer rows.Close()

	var blacklistedVoters []database.BlacklistType

	for rows.Next() {
		var voter database.BlacklistType

		err := rows.Scan(&voter.Address, &voter.Duration, &voter.IsActive, &voter.CreatedAt)

		if err != nil {
			log.Println(err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan blacklisted voter"})
			return
		}

		blacklistedVoters = append(blacklistedVoters, voter)
	}

	c.JSON(http.StatusOK, gin.H{"blacklist": blacklistedVoters})
}

func GetIsVoterBlacklistedHandler(c *gin.Context) {
	address := c.Param("address")

	rows, err := database.Get().Query(
		c.Request.Context(),
		`SELECT address, (EXTRACT(EPOCH FROM duration) / 60)::int, is_active, created_at
		FROM blacklist WHERE address = $1 AND is_active = TRUE`,
		address,
	)

	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get blacklisted voter"})
		return
	}

	var blacklist []database.BlacklistType

	for rows.Next() {
		var voter database.BlacklistType
		err := rows.Scan(&voter.Address, &voter.Duration, &voter.IsActive, &voter.CreatedAt)

		if err != nil {
			log.Println(err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan blacklisted voter"})
			return
		}

		blacklist = append(blacklist, voter)
	}

	c.JSON(http.StatusOK, gin.H{"blacklist": blacklist})
}

func RemoveVoterFromBlacklistHandler(c *gin.Context) {
	// TODO: implement this
}
