package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// CreateReservationHandler registra uma nova reserva.
// @Summary Cria reserva
// @Description Cria uma nova reserva para um espaço.
// @Tags reservas
// @Accept json
// @Produce json
// @Param input body dto.CreateReservationDTO true "Dados da reserva"
// @Success 201 {object} dto.ReservationOutputDTO
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /reservas [post]
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

// GetReservationByIDHandler retorna uma reserva pelo ID.
// @Summary Busca reserva por ID
// @Description Retorna os dados de uma reserva específica.
// @Tags reservas
// @Produce json
// @Param id path int true "ID da reserva"
// @Success 200 {object} dto.ReservationOutputDTO
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /reservas/{id} [get]
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

// GetAllReservationsHandler lista todas as reservas.
// @Summary Lista reservas
// @Description Lista todas as reservas cadastradas.
// @Tags reservas
// @Produce json
// @Success 200 {array} dto.ReservationOutputDTO
// @Failure 500 {object} gin.H
// @Router /reservas [get]
func GetAllReservationsHandler(c *gin.Context) {
	reservations, err := services.GetAllReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar reservas"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}
