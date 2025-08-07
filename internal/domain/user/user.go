package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	CPF      string `json:"cpf"`
	Password string `json:"_"`
	Role     string `json:"role"`
}

type Repository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByCPF(cpf string) (*User, error)
	FindAll() ([]User, error)
	DeleteByID(id uint) error
}
