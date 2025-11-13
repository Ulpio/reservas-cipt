package dto

import "time"

// ClienteInputDTO usado para criar/atualizar cliente
type ClienteInputDTO struct {
	Name      string `json:"name" binding:"required,min=3,max=255"`
	CPF       string `json:"cpf" binding:"required"`
	BirthDate string `json:"birth_date" binding:"required"` // Formato: YYYY-MM-DD
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// ClienteOutputDTO retornado nas respostas da API
type ClienteOutputDTO struct {
	ID                 uint      `json:"id"`
	ClientID           uint      `json:"client_id"` // Alias para compatibilidade
	Name               string    `json:"name"`
	CPF                string    `json:"cpf"`
	BirthDate          string    `json:"birth_date"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Strikes            int       `json:"strikes"`
	StrikesCount       int       `json:"strikes_count"`       // Alias para compatibilidade
	TotalReservations  int       `json:"total_reservations"`  // Quantidade de reservas
	CreatedAt          time.Time `json:"created_at"`
}
