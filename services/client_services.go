package services

import (
	"errors"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/utils"
	"gorm.io/gorm"
)

func BuscarOuCriarCliente(input dto.ClienteInputDTO) (dto.ClienteOutputDTO, error) {
	// Normalizar CPF e telefone
	normalizedCPF := utils.NormalizeCPF(input.CPF)
	normalizedPhone := utils.NormalizePhone(input.Phone)

	var cliente models.Client
	result := database.DB.Where("cpf = ?", normalizedCPF).First(&cliente)

	if result.Error == nil {
		return toClientOutput(cliente), nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return dto.ClienteOutputDTO{}, result.Error
	}

	// Parsear data de nascimento
	birthDate, err := utils.ParseDate(input.BirthDate)
	if err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	novo := models.Client{
		Name:      input.Name,
		CPF:       normalizedCPF,
		BirthDate: birthDate,
		Email:     input.Email,
		Phone:     normalizedPhone,
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
	normalizedCPF := utils.NormalizeCPF(cpf)
	var cliente models.Client
	if err := database.DB.Where("cpf = ?", normalizedCPF).First(&cliente).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}
	return toClientOutput(cliente), nil
}

// GetClientByPhone busca um cliente pelo telefone
func GetClientByPhone(phone string) (dto.ClienteOutputDTO, error) {
	normalizedPhone := utils.NormalizePhone(phone)
	var cliente models.Client
	if err := database.DB.Where("phone = ?", normalizedPhone).First(&cliente).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}
	return toClientOutput(cliente), nil
}

// GetClientByID busca um cliente pelo ID
func GetClientByID(id uint) (dto.ClienteOutputDTO, error) {
	var cliente models.Client
	if err := database.DB.First(&cliente, id).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}
	return toClientOutput(cliente), nil
}

// CreateClient cria um novo cliente com todas as validações
func CreateClient(input dto.ClienteInputDTO) (dto.ClienteOutputDTO, error) {
	// 1. Validar CPF
	if err := utils.ValidateCPF(input.CPF); err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	// 2. Validar email (se fornecido)
	if err := utils.ValidateEmail(input.Email); err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	// 3. Validar telefone (se fornecido)
	if err := utils.ValidatePhoneNumber(input.Phone); err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	// 4. Parsear e validar data de nascimento
	birthDate, err := utils.ParseDate(input.BirthDate)
	if err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	// Validar se não é data futura
	if err := utils.ValidateFutureDate(birthDate); err == nil {
		// Se ValidateFutureDate retornou nil, significa que é hoje ou futuro (erro!)
		return dto.ClienteOutputDTO{}, errors.New("data de nascimento deve ser uma data passada")
	}

	// 5. Validar idade mínima (18 anos)
	if err := utils.ValidateMinimumAge(birthDate, 18); err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	// 6. Normalizar CPF e telefone
	normalizedCPF := utils.NormalizeCPF(input.CPF)
	normalizedPhone := utils.NormalizePhone(input.Phone)

	// 7. Verificar se CPF já existe
	var existingByCPF models.Client
	if err := database.DB.Where("cpf = ?", normalizedCPF).First(&existingByCPF).Error; err == nil {
		return dto.ClienteOutputDTO{}, errors.New("CPF já cadastrado no sistema")
	}

	// 8. Verificar se telefone já existe (apenas se fornecido)
	if normalizedPhone != "" {
		var existingByPhone models.Client
		if err := database.DB.Where("phone = ?", normalizedPhone).First(&existingByPhone).Error; err == nil {
			return dto.ClienteOutputDTO{}, errors.New("Telefone já cadastrado")
		}
	}

	// 9. Criar novo cliente
	novo := models.Client{
		Name:      input.Name,
		CPF:       normalizedCPF,
		BirthDate: birthDate,
		Email:     input.Email,
		Phone:     normalizedPhone,
	}

	if err := database.DB.Create(&novo).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}

	return toClientOutput(novo), nil
}

func UpdateClient(id uint, input dto.ClienteInputDTO) (dto.ClienteOutputDTO, error) {
	var cliente models.Client
	if err := database.DB.First(&cliente, id).Error; err != nil {
		return dto.ClienteOutputDTO{}, err
	}
	cliente.Name = input.Name
	cliente.Email = input.Email
	cliente.Phone = input.Phone

	if err := database.DB.Save(&cliente).Error; err != nil {
		return dto.ClienteOutputDTO{}, nil
	}
	return toClientOutput(cliente), nil
}

func toClientOutput(client models.Client) dto.ClienteOutputDTO {
	// Contar total de reservas
	var totalReservations int64
	database.DB.Model(&models.Reservation{}).Where("client_id = ?", client.ID).Count(&totalReservations)

	return dto.ClienteOutputDTO{
		ID:                client.ID,
		ClientID:          client.ID, // Alias para compatibilidade
		Name:              client.Name,
		CPF:               client.CPF,
		BirthDate:         client.BirthDate.Format("2006-01-02"),
		Email:             client.Email,
		Phone:             client.Phone,
		Strikes:           client.Strikes,
		StrikesCount:      client.Strikes, // Alias para compatibilidade
		TotalReservations: int(totalReservations),
		CreatedAt:         client.CreatedAt,
	}
}
