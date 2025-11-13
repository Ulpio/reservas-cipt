package dto

import "time"

// ReservaStatusDTO representa uma reserva simplificada para o status
type ReservaStatusDTO struct {
	ID         uint      `json:"id"`
	Cliente    string    `json:"cliente"`
	HoraInicio time.Time `json:"horaInicio"`
	HoraFim    time.Time `json:"horaFim"`
}

// SpaceStatusDTO representa o status completo de um espaço
type SpaceStatusDTO struct {
	ID            uint                `json:"id"`
	Nome          string              `json:"nome"`
	Status        string              `json:"status"` // "disponivel", "ocupado", "manutencao"
	ReservaAtual  *ReservaStatusDTO   `json:"reservaAtual"`
	ProximaReserva *ReservaStatusDTO  `json:"proximaReserva"`
	ReservasHoje  []ReservaStatusDTO  `json:"reservasHoje"`
}

