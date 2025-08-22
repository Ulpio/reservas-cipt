package services

import (
	"errors"

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

func GetAllClientes() ([]dto.ClienteOutputDTO, error) {
	var clientes []models.Client
	if err := database.DB.Find(&clientes); err != nil {
		return nil, err.Error
	}
	var output []dto.ClienteOutputDTO
	for _, s := range clientes {
		output = append(output, toClientOutput(s))
	}
	return output, nil
}

func GetClientByCPF(cpf string) (dto.ClienteOutputDTO, error) {
	var cliente models.Client
	if err := database.DB.Where("cpf = ? ", cpf).Find(&cliente); err != nil {
		return dto.ClienteOutputDTO{}, err.Error
	}
	if cliente.ID == 0 {
		return dto.ClienteOutputDTO{}, errors.New("erro ao encontrar usuario")
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
