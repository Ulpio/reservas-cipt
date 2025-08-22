package dto

import "time"

// CreateReservationDTO contains data needed to create a reservation.
type CreateReservationDTO struct {
	ClientID       uint      `json:"client_id" binding:"required"`
	ReceptionistID uint      `json:"receptionist_id" binding:"required"`
	SpaceID        uint      `json:"space_id" binding:"required"`
	Date           time.Time `json:"date" binding:"required"`
	StartTime      time.Time `json:"start_time" binding:"required"`
	DurationHours  int       `json:"duration_hours" binding:"required"`
}

// ReservationOutputDTO represents reservation data returned in responses.
type ReservationOutputDTO struct {
	ID             uint      `json:"id"`
	ClientID       uint      `json:"client_id"`
	ReceptionistID uint      `json:"receptionist_id"`
	SpaceID        uint      `json:"space_id"`
	Date           time.Time `json:"date"`
	StartTime      time.Time `json:"start_time"`
	DurationHours  int       `json:"duration_hours"`
}
