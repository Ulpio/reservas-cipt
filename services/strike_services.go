package services

import (
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

func CreateStrike(input dto.StrikeInputDTO) (dto.StrikeOutputDTO, error) {
	strike := models.Strike{
		ClientID: input.ClientID,
		Reason:   input.Reason,
		Photo:    input.Photo,
	}
	if err := database.DB.Create(&strike).Error; err != nil {
		return dto.StrikeOutputDTO{}, err
	}
	if err := updateClientStrikeCount(strike.ClientID); err != nil {
		return dto.StrikeOutputDTO{}, err
	}
	return toStrikeOutput(strike), nil
}

func GetStrikesByClient(clientID uint) ([]dto.StrikeOutputDTO, error) {
	var strikes []models.Strike
	if err := database.DB.Where("client_id = ?", clientID).Find(&strikes).Error; err != nil {
		return nil, err
	}
	var output []dto.StrikeOutputDTO
	for _, s := range strikes {
		output = append(output, toStrikeOutput(s))
	}
	return output, nil
}

func RevokeStrike(id uint) (dto.StrikeOutputDTO, error) {
	var strike models.Strike
	if err := database.DB.First(&strike, id).Error; err != nil {
		return dto.StrikeOutputDTO{}, err
	}
	if strike.Revoked {
		return toStrikeOutput(strike), nil
	}
	now := time.Now()
	strike.Revoked = true
	strike.RevokedAt = &now
	if err := database.DB.Save(&strike).Error; err != nil {
		return dto.StrikeOutputDTO{}, err
	}
	if err := updateClientStrikeCount(strike.ClientID); err != nil {
		return dto.StrikeOutputDTO{}, err
	}
	return toStrikeOutput(strike), nil
}

func toStrikeOutput(s models.Strike) dto.StrikeOutputDTO {
	return dto.StrikeOutputDTO{
		ID:        s.ID,
		ClientID:  s.ClientID,
		Reason:    s.Reason,
		Photo:     s.Photo,
		CreatedAt: s.CreatedAt,
		Revoked:   s.Revoked,
		RevokedAt: s.RevokedAt,
	}
}

func updateClientStrikeCount(clientID uint) error {
	var count int64
	if err := database.DB.Model(&models.Strike{}).Where("client_id = ? AND revoked = ?", clientID, false).Count(&count).Error; err != nil {
		return err
	}
	return database.DB.Model(&models.Client{}).Where("id = ?", clientID).Update("strikes", int(count)).Error
}
