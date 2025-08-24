package dto

import "time"

type StrikeInputDTO struct {
	ClientID uint   `json:"client_id" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
	Photo    string `json:"photo"`
}

type StrikeOutputDTO struct {
	ID        uint       `json:"id"`
	ClientID  uint       `json:"client_id"`
	Reason    string     `json:"reason"`
	Photo     string     `json:"photo"`
	CreatedAt time.Time  `json:"created_at"`
	Revoked   bool       `json:"revoked"`
	RevokedAt *time.Time `json:"revoked_at"`
}
