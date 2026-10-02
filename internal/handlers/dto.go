package handlers

type RegisterVoterDTO struct {
	Id      string `json:"id" binding:"required,len=10,numeric"`
	Address string `json:"address" binding:"required,eth_addr"`
}
