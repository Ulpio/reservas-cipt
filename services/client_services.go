package services

import (
	"errors"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
	"gorm.io/gorm"
)

func BuscarOuCriarCliente(input dto.ClienteInputDTO) (dto.ClienteOutputDTO, error) {
	var cliente models.Client
	result := database.DB.Where("cpf =?", input.CPF).First(&cliente)

	if result.Error == nil {
		return toClientOutput(cliente), nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dto.ClienteOutputDTO{}, result.Error
	}

	novo := models.Client{
		Name:  input.Name,
		CPF:   input.CPF,
		Email: input.Email,
		Phone: input.Phone,
	}
	if err := database.DB.Create(&novo).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}
	return toClientOutput(novo), nil
}

func GetAllClientes() ([]dto.ClienteOutputDTO, error) {
	var clientes []models.Client
	if err := database.DB.Find(&clientes).Error; err != nil {
		return nil, err
	}
	var output []dto.ClienteOutputDTO
	for _, s := range clientes {
		output = append(output, toClientOutput(s))
	}
	return output, nil
}

func GetClientByCPF(cpf string) (dto.ClienteOutputDTO, error) {
	var cliente models.Client
	if err := database.DB.Where("cpf = ?", cpf).First(&cliente).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}
	return toClientOutput(cliente), nil
}

func toClientOutput(client models.Client) dto.ClienteOutputDTO {
	return dto.ClienteOutputDTO{
		ID:      client.ID,
		Name:    client.Name,
		CPF:     client.CPF,
		Email:   client.Email,
		Phone:   client.Phone,
		Strikes: client.Strikes,
	}
}
