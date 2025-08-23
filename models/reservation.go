package models

import "time"

// Reservation represents a booking of a space by a client handled by a receptionist
// Data includes client, receptionist, space, date, start time and duration in hours.
type Reservation struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ClientID       uint      `gorm:"not null"`
	ReceptionistID uint      `gorm:"not null"`
	SpaceID        uint      `gorm:"not null"`
	Date           time.Time `gorm:"not null"`
	StartTime      time.Time `gorm:"not null"`
	DurationHours  int       `gorm:"not null"`
}
