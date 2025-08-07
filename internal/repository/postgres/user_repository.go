package postgres

import (
	"github.com/Ulpio/reservas-cipt/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) Create(user *user.User) error {
	return u.DB.Create(user).Error
}

func (u *UserRepository) FindByID(id uint) (*user.User, error) {
	var usr user.User
	if err := u.DB.First(&usr, id).Error; err != nil {
		return nil, err
	}
	return &usr, nil
}

func (u *UserRepository) FindByCPF(cpf string) (*user.User, error) {
	var usr user.User
	if err := u.DB.Where("cpf = ?", cpf).First(&usr).Error; err != nil {
		return nil, err
	}
	return &usr, nil
}

func (u *UserRepository) FindAll() ([]user.User, error) {
	var users []user.User
	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserRepository) DeleteByID(id uint) error {
	return u.DB.Delete(&user.User{}, id).Error
}
