package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name            string
	Password        string
	CPF             string `gorm:"unique"`
	Role            string
	DefaultPassword bool
}
