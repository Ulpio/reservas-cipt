package services

import (
	"errors"
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

// ValidateSpaceType valida se o tipo de espaço é válido
func ValidateSpaceType(spaceType string) error {
	validTypes := []string{"visitantes", "permissionarios"}
	for _, valid := range validTypes {
		if spaceType == valid {
			return nil
		}
	}
	return errors.New("tipo de espaço inválido. Use 'visitantes' ou 'permissionarios'")
}

func CreateSpace(input dto.CreateSpaceDTO) (dto.SpaceOutputDTO, error) {
	// Validar tipo de espaço
	if err := ValidateSpaceType(input.Type); err != nil {
		return dto.SpaceOutputDTO{}, err
	}

	space := models.Space{
		Name:     input.Name,
		Type:     input.Type,
		Status:   input.Status,
		Notice:   input.Notice,
		Capacity: uint(input.Capacity),
	}

	if err := database.DB.Create(&space).Error; err != nil {
		return dto.SpaceOutputDTO{}, err
	}

	return toSpaceOutput(space), nil
}

func GetAllSpaces() ([]dto.SpaceOutputDTO, error) {
	var spaces []models.Space
	if err := database.DB.Find(&spaces).Error; err != nil {
		return nil, err
	}
	var output []dto.SpaceOutputDTO
	for _, s := range spaces {
		output = append(output, toSpaceOutputWithRealTimeStatus(s))
	}
	return output, nil
}

func GetSpaceByID(id uint) (dto.SpaceOutputDTO, error) {
	var space models.Space
	if err := database.DB.First(&space, id).Error; err != nil {
		return dto.SpaceOutputDTO{}, nil
	}
	return toSpaceOutput(space), nil
}

func UpdateSpace(id uint, input dto.UpdateSpaceDTO) (dto.SpaceOutputDTO, error) {
	var space models.Space
	if err := database.DB.First(&space, id).Error; err != nil {
		return dto.SpaceOutputDTO{}, err
	}

	// Validar tipo de espaço se fornecido
	if input.Type != "" {
		if err := ValidateSpaceType(input.Type); err != nil {
			return dto.SpaceOutputDTO{}, err
		}
		space.Type = input.Type
	}

	space.Name = input.Name
	space.Status = input.Status
	space.Notice = input.Notice
	space.Capacity = uint(input.Capacity)

	if err := database.DB.Save(&space).Error; err != nil {
		return dto.SpaceOutputDTO{}, err
	}
	return toSpaceOutput(space), nil
}

func DeleteSpace(id uint) error {
	return database.DB.Delete(&models.Space{}, id).Error
}

func UpdateSpaceStatus(id uint, status string) error {
	return database.DB.Model(&models.Space{}).Where("id = ?", id).Update("status", status).Error
}

func UpdateSpaceNotice(id uint, notice string) error {
	return database.DB.Model(&models.Space{}).Where("id = ?", id).Update("notice", notice).Error
}

func toSpaceOutput(space models.Space) dto.SpaceOutputDTO {
	return dto.SpaceOutputDTO{
		ID:       space.ID,
		Name:     space.Name,
		Type:     space.Type,
		Status:   space.Status,
		Notice:   space.Notice,
		Capacity: int(space.Capacity),
	}
}

// toSpaceOutputWithRealTimeStatus retorna o espaço com status calculado em tempo real
func toSpaceOutputWithRealTimeStatus(space models.Space) dto.SpaceOutputDTO {
	// Se o espaço está em manutenção, manter esse status
	if space.Status == "manutencao" || space.Status == "manutenção" {
		return dto.SpaceOutputDTO{
			ID:       space.ID,
			Name:     space.Name,
			Type:     space.Type,
			Status:   "manutencao",
			Notice:   space.Notice,
			Capacity: int(space.Capacity),
		}
	}

	// Verificar se há reserva ativa agora
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := today.Add(24 * time.Hour)

	var activeReservations []models.Reservation
	database.DB.Where("space_id = ? AND date >= ? AND date < ?", space.ID, today, endOfDay).Find(&activeReservations)

	// Verificar se alguma reserva está ativa agora
	for _, reservation := range activeReservations {
		endTime := reservation.StartTime.Add(time.Duration(reservation.DurationHours) * time.Hour)
		
		// Se a hora atual está entre start_time e end_time
		if now.After(reservation.StartTime) && now.Before(endTime) {
			return dto.SpaceOutputDTO{
				ID:       space.ID,
				Name:     space.Name,
				Type:     space.Type,
				Status:   "ocupado",
				Notice:   space.Notice,
				Capacity: int(space.Capacity),
			}
		}
	}

	// Se não está ocupado e não está em manutenção, está disponível
	return dto.SpaceOutputDTO{
		ID:       space.ID,
		Name:     space.Name,
		Type:     space.Type,
		Status:   "ativo",
		Notice:   space.Notice,
		Capacity: int(space.Capacity),
	}
}
