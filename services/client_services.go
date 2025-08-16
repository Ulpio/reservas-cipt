package services

import (
	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

func BuscarOuCriarCliente(input dto.ClienteInputDTO) (dto.ClienteOutputDTO, error) {
	var cliente models.Client
	result := database.DB.Where("cpf =?", input.CPF).First(&cliente)

	if result.Error == nil {
		return dto.ClienteOutputDTO{
			ID:      cliente.ID,
			Name:    cliente.Name,
			CPF:     cliente.CPF,
			Email:   cliente.Email,
			Phone:   cliente.Phone,
			Strikes: cliente.Strikes,
		}, nil
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
	return dto.ClienteOutputDTO{
		ID:      novo.ID,
		Name:    novo.Name,
		CPF:     novo.CPF,
		Email:   novo.Email,
		Phone:   novo.Phone,
		Strikes: novo.Strikes,
	}, nil
}
