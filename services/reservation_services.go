package services

import (
	"time"

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

	output, err := toReservationOutput(reservation)
	if err != nil {
		return dto.ReservationOutputDTO{}, err
	}
	return output, nil
}

// GetReservationByID retrieves a reservation by its ID.
func GetReservationByID(id uint) (dto.ReservationOutputDTO, error) {
	var reservation models.Reservation
	if err := database.DB.First(&reservation, id).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}
	return toReservationOutput(reservation)
}

// GetAllReservations returns all reservations.
func GetAllReservations() ([]dto.ReservationOutputDTO, error) {
	var reservations []models.Reservation
	if err := database.DB.Find(&reservations).Error; err != nil {
		return nil, err
	}
	var output []dto.ReservationOutputDTO
	for _, r := range reservations {
		o, err := toReservationOutput(r)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}

func toReservationOutput(r models.Reservation) (dto.ReservationOutputDTO, error) {
	var client models.Client
	if err := database.DB.First(&client, r.ClientID).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}
	var receptionist models.User
	if err := database.DB.First(&receptionist, r.ReceptionistID).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}
	var space models.Space
	if err := database.DB.First(&space, r.SpaceID).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	end := r.StartTime.Add(time.Duration(r.DurationHours) * time.Hour)

	return dto.ReservationOutputDTO{
		ID:               r.ID,
		ClientName:       client.Name,
		ReceptionistName: receptionist.Name,
		SpaceName:        space.Name,
		Date:             r.Date.Format("02/01/2006"),
		StartTime:        r.StartTime.Format("15:04"),
		DurationHours:    r.DurationHours,
		EndTime:          end.Format("15:04"),
	}, nil
}
