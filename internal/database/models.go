package database

import "time"

type VoterType struct {
	Id        string    `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type BlacklistType struct {
	Address   string    `json:"address"`
	Duration  int       `json:"duration"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
