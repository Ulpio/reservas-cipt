package services

import (
	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

// CreateReservation registers a new reservation of a space for a client handled by a receptionist.
func CreateReservation(input dto.CreateReservationDTO) (dto.ReservationOutputDTO, error) {
	reservation := models.Reservation{
		ClientID:       input.ClientID,
		ReceptionistID: input.ReceptionistID,
		SpaceID:        input.SpaceID,
		Date:           input.Date,
		StartTime:      input.StartTime,
		DurationHours:  input.DurationHours,
	}

	if err := database.DB.Create(&reservation).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	return toReservationOutput(reservation), nil
}

// GetReservationByID retrieves a reservation by its ID.
func GetReservationByID(id uint) (dto.ReservationOutputDTO, error) {
	var reservation models.Reservation
	if err := database.DB.First(&reservation, id).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}
	return toReservationOutput(reservation), nil
}

// GetAllReservations returns all reservations.
func GetAllReservations() ([]dto.ReservationOutputDTO, error) {
	var reservations []models.Reservation
	if err := database.DB.Find(&reservations).Error; err != nil {
		return nil, err
	}
	var output []dto.ReservationOutputDTO
	for _, r := range reservations {
		output = append(output, toReservationOutput(r))
	}
	return output, nil
}

func toReservationOutput(r models.Reservation) dto.ReservationOutputDTO {
	return dto.ReservationOutputDTO{
		ID:             r.ID,
		ClientID:       r.ClientID,
		ReceptionistID: r.ReceptionistID,
		SpaceID:        r.SpaceID,
		Date:           r.Date,
		StartTime:      r.StartTime,
		DurationHours:  r.DurationHours,
	}
}
