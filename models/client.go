package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Name  string `gorm:"not null"`
	CPF   string `gorm:"unique;not null"`
	Email string
	Phone string

	Strikes int
}
