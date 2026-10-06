package handlers

type RegisterVoterDTO struct {
	ID      string `json:"id" binding:"required,len=10,numeric"`
	Address string `json:"address" binding:"required,eth_addr"`
}

type BlacklistVoterDTO struct {
	Address  string `json:"address" binding:"required,eth_addr"`
	Duration int    `json:"duration" binding:"required,min=1,max=1440"`
}
