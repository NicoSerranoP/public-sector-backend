package handlers

import "github.com/gin-gonic/gin"

func RelayCastVoteHandler(c *gin.Context) {
	// TODO: save tx data to send in batches every N seconds

	// TODO: check `to` is the VOTING_PLUGIN_ADDRESS
	// TODO: check `selector` is vote()
	// TODO: check `from` is not in blacklist
	// TODO: check plugin.canVote() returns true
}

func RelayCreateProposalHandler(c *gin.Context) {}
