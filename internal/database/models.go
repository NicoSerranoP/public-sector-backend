package database

import "time"

type VoterType struct {
	Id        string    `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
