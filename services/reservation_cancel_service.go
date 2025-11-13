package services

import (
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

// CancelReservation cancela uma reserva e envia email de notificação
func CancelReservation(id uint) error {
	// 1. Buscar a reserva
	var reservation models.Reservation
	if err := database.DB.First(&reservation, id).Error; err != nil {
		return err
	}

	// 2. Buscar dados relacionados para o email
	var client models.Client
	var space models.Space
	database.DB.First(&client, reservation.ClientID)
	database.DB.First(&space, reservation.SpaceID)

	startTime := reservation.StartTime.Format("15:04")
	endTime := reservation.StartTime.Add(time.Duration(reservation.DurationHours) * time.Hour).Format("15:04")

	// 3. Deletar a reserva
	if err := database.DB.Delete(&reservation).Error; err != nil {
		return err
	}

	// 4. Enviar email de cancelamento (assíncrono)
	go sendReservationCancellationEmail(client, space, reservation, startTime, endTime)

	return nil
}

// GetReservationDetails retorna os detalhes de uma reserva (para validações)
func GetReservationDetails(id uint) (dto.ReservationOutputDTO, error) {
	return GetReservationByID(id)
}

