package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// CreateReservationHandler handles the creation of a reservation.
func CreateReservationHandler(c *gin.Context) {
	var input dto.CreateReservationDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	reservation, err := services.CreateReservation(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar reserva"})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// GetReservationByIDHandler returns a reservation by its ID.
func GetReservationByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reservation, err := services.GetReservationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// GetAllReservationsHandler lists all reservations.
func GetAllReservationsHandler(c *gin.Context) {
	reservations, err := services.GetAllReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar reservas"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}
