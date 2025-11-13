package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name      string    `gorm:"not null"`
	CPF       string    `gorm:"unique;not null"`
	BirthDate time.Time `gorm:"not null"`
	Email     string
	Phone     string

	Strikes int
}
