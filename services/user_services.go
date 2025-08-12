package services

import (
	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/utils"
)

func GetAllUsers() ([]dto.UserOutputDTO, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	var output []dto.UserOutputDTO
	for _, u := range users {
		output = append(output, dto.UserOutputDTO{
			ID:   u.ID,
			Name: u.Name,
			CPF:  u.CPF,
			Role: u.Role,
		})
	}
	return output, nil
}

func GetUserByID(id uint) (dto.UserOutputDTO, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return dto.UserOutputDTO{}, err
	}
	return dto.UserOutputDTO{
		ID:   user.ID,
		Name: user.Name,
		CPF:  user.CPF,
		Role: user.Role,
	}, nil
}

func CreateUser(name, cpf, role string) (dto.UserOutputDTO, error) {
	hashedPass, err := utils.HashPassword("12345678")
	if err != nil {
		return dto.UserOutputDTO{}, err
	}
	user := models.User{
		Name: name, CPF: cpf, Role: role, Password: hashedPass,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return dto.UserOutputDTO{}, err
	}
	output := dto.UserOutputDTO{
		ID: user.ID, CPF: user.CPF, Role: user.Role, Name: user.Name,
	}
	return output, nil
}

func DeleteUser(id uint) error {
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
