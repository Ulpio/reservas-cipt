package services

import (
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

// GetAllSpacesStatus retorna o status de todos os espaços com suas reservas
func GetAllSpacesStatus() ([]dto.SpaceStatusDTO, error) {
	var spaces []models.Space
	if err := database.DB.Find(&spaces).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var result []dto.SpaceStatusDTO

	for _, space := range spaces {
		spaceStatus := dto.SpaceStatusDTO{
			ID:            space.ID,
			Nome:          space.Name,
			Status:        determineSpaceStatus(space, now),
			ReservaAtual:  nil,
			ProximaReserva: nil,
			ReservasHoje:  []dto.ReservaStatusDTO{},
		}

		// Se o espaço está em manutenção, não precisa buscar reservas
		if space.Status == "manutencao" {
			result = append(result, spaceStatus)
			continue
		}

		// Buscar todas as reservas do dia para este espaço com JOIN para pegar o nome do cliente
		type ReservationWithClient struct {
			ReservationID uint
			ClientName    string
			StartTime     time.Time
			DurationHours int
		}

		var reservationsWithClients []ReservationWithClient
		if err := database.DB.Table("reservations").
			Select("reservations.id as reservation_id, clients.name as client_name, reservations.start_time, reservations.duration_hours").
			Joins("JOIN clients ON clients.id = reservations.client_id").
			Where("reservations.space_id = ? AND reservations.date >= ? AND reservations.date < ?", space.ID, startOfDay, endOfDay).
			Order("reservations.start_time ASC").
			Scan(&reservationsWithClients).Error; err != nil {
			return nil, err
		}

		// Processar reservas
		var reservasHoje []dto.ReservaStatusDTO
		var reservaAtual *dto.ReservaStatusDTO
		var proximaReserva *dto.ReservaStatusDTO

		for _, res := range reservationsWithClients {
			reservaDTO := dto.ReservaStatusDTO{
				ID:         res.ReservationID,
				Cliente:    res.ClientName,
				HoraInicio: res.StartTime,
				HoraFim:    calculateEndTime(res.StartTime, res.DurationHours),
			}

			reservasHoje = append(reservasHoje, reservaDTO)

			endTime := calculateEndTime(res.StartTime, res.DurationHours)

			// Verificar se é a reserva atual (em andamento)
			if now.After(res.StartTime) && now.Before(endTime) {
				reservaAtualCopy := reservaDTO
				reservaAtual = &reservaAtualCopy
			}

			// Verificar se é a próxima reserva
			if proximaReserva == nil && res.StartTime.After(now) {
				proximaReservaCopy := reservaDTO
				proximaReserva = &proximaReservaCopy
			}
		}

		spaceStatus.ReservasHoje = reservasHoje
		spaceStatus.ReservaAtual = reservaAtual
		spaceStatus.ProximaReserva = proximaReserva

		// Se tem reserva atual, status deve ser "ocupado"
		if reservaAtual != nil {
			spaceStatus.Status = "ocupado"
		}

		result = append(result, spaceStatus)
	}

	return result, nil
}

// determineSpaceStatus determina o status base do espaço
func determineSpaceStatus(space models.Space, now time.Time) string {
	if space.Status == "manutencao" {
		return "manutencao"
	}
	// Status inicial é disponível, pode ser alterado para ocupado se houver reserva atual
	return "disponivel"
}

// calculateEndTime calcula o horário de término da reserva
func calculateEndTime(startTime time.Time, durationHours int) time.Time {
	return startTime.Add(time.Duration(durationHours) * time.Hour)
}

