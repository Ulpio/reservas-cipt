package services

import (
	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

func CreateSpace(input dto.CreateSpaceDTO) (dto.SpaceOutputDTO, error) {
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
		output = append(output, toSpaceOutput(s))
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
	space.Name = input.Name
	space.Type = input.Type
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
