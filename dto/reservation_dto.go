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
	ID               uint   `json:"id"`
	ClientName       string `json:"client_name"`
	ReceptionistName string `json:"receptionist_name"`
	SpaceName        string `json:"space_name"`
	Date             string `json:"date"`
	StartTime        string `json:"start_time"`
	DurationHours    int    `json:"duration_hours"`
	EndTime          string `json:"end_time"`
}
