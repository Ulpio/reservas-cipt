package services

import (
	"errors"
	"log"
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/utils"
)

var emailService *EmailService

func init() {
	emailService = NewEmailService()
}

// CreateReservationFromInput cria uma reserva a partir do input com strings
func CreateReservationFromInput(input dto.CreateReservationInputDTO) (dto.ReservationOutputDTO, error) {
	// 1. Parsear e validar data
	date, err := utils.ParseDate(input.Date)
	if err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	// Validar que a data não é passada
	if err := utils.ValidateFutureDate(date); err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	// 2. Parsear e validar hora
	startTime, err := utils.ParseDateTime(input.Date, input.StartTime)
	if err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	// 3. Validar que cliente existe
	var client models.Client
	if err := database.DB.First(&client, input.ClientID).Error; err != nil {
		return dto.ReservationOutputDTO{}, errors.New("cliente não encontrado")
	}

	// 4. Validar que recepcionista existe e tem permissão
	var receptionist models.User
	if err := database.DB.First(&receptionist, input.ReceptionistID).Error; err != nil {
		return dto.ReservationOutputDTO{}, errors.New("recepcionista não encontrado")
	}

	// Verificar se o usuário tem role adequada
	role := receptionist.Role
	if role != "admin" && role != "recepcionista" && role != "locacao" {
		return dto.ReservationOutputDTO{}, errors.New("usuário não tem permissão para criar reservas")
	}

	// 5. Validar que espaço existe e está ativo
	var space models.Space
	if err := database.DB.First(&space, input.SpaceID).Error; err != nil {
		return dto.ReservationOutputDTO{}, errors.New("espaço não encontrado")
	}

	if space.Status != "ativo" {
		return dto.ReservationOutputDTO{}, errors.New("espaço não está disponível para reservas")
	}

	// 6. Validar duration_hours
	if input.DurationHours < 1 || input.DurationHours > 24 {
		return dto.ReservationOutputDTO{}, errors.New("duração deve ser entre 1 e 24 horas")
	}

	// 7. Validar conflito de horários
	// Buscar todas as reservas do espaço no mesmo dia
	var existingReservations []models.Reservation
	if err := database.DB.Where("space_id = ? AND date = ?", input.SpaceID, date).Find(&existingReservations).Error; err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	// Verificar sobreposição de horários
	endTime := startTime.Add(time.Duration(input.DurationHours) * time.Hour)
	for _, existing := range existingReservations {
		existingEnd := existing.StartTime.Add(time.Duration(existing.DurationHours) * time.Hour)
		
		// Verificar se há sobreposição
		// Nova reserva começa antes do fim da existente E termina depois do início da existente
		if startTime.Before(existingEnd) && endTime.After(existing.StartTime) {
			return dto.ReservationOutputDTO{}, errors.New("espaço já reservado para este horário")
		}
	}

	// 8. Criar a reserva
	internalDTO := dto.CreateReservationDTO{
		ClientID:       input.ClientID,
		ReceptionistID: input.ReceptionistID,
		SpaceID:        input.SpaceID,
		Date:           date,
		StartTime:      startTime,
		DurationHours:  input.DurationHours,
	}

	reservation, err := CreateReservation(internalDTO)
	if err != nil {
		return dto.ReservationOutputDTO{}, err
	}

	// 9. Enviar email de confirmação (assíncrono)
	go sendReservationConfirmationEmail(client, space, models.Reservation{
		ClientID:       input.ClientID,
		ReceptionistID: input.ReceptionistID,
		SpaceID:        input.SpaceID,
		Date:           date,
		StartTime:      startTime,
		DurationHours:  input.DurationHours,
	}, startTime.Format("15:04"), endTime.Format("15:04"))

	return reservation, nil
}

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

// sendReservationConfirmationEmail envia email de confirmação (goroutine)
func sendReservationConfirmationEmail(client models.Client, space models.Space, reservation models.Reservation, startTime, endTime string) {
	if err := emailService.SendReservationConfirmation(client, space, reservation, startTime, endTime); err != nil {
		log.Printf("⚠️  Erro ao enviar email de confirmação: %v\n", err)
	}
}

// sendReservationCancellationEmail envia email de cancelamento (goroutine)
func sendReservationCancellationEmail(client models.Client, space models.Space, reservation models.Reservation, startTime, endTime string) {
	if err := emailService.SendReservationCancellation(client, space, reservation, startTime, endTime); err != nil {
		log.Printf("⚠️  Erro ao enviar email de cancelamento: %v\n", err)
	}
}
