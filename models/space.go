package models

import "gorm.io/gorm"

type Space struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Type     string `gorm:"not null"` // Tipos: visitantes, permissionarios (reserva de salas de reunião)
	Status   string `gorm:"default:ativo"`
	Notice   string
	Capacity uint `gorm:"not null"`
}
