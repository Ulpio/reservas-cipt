package models

import (
	"time"

	"gorm.io/gorm"
)

type Strike struct {
	gorm.Model
	ClientID  uint
	Reason    string `gorm:"not null"`
	Photo     string
	Revoked   bool
	RevokedAt *time.Time
}
