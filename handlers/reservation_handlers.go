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
// @Description Cria uma nova reserva para um espaço com validações completas.
// @Tags reservas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.CreateReservationInputDTO true "Dados da reserva"
// @Success 201 {object} dto.ReservationOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservas [post]
func CreateReservationHandler(c *gin.Context) {
	var input dto.CreateReservationInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	reservation, err := services.CreateReservationFromInput(input)
	if err != nil {
		errMsg := err.Error()
		
		// Retornar status code apropriado baseado no erro
		switch errMsg {
		case "cliente não encontrado", "recepcionista não encontrado", "espaço não encontrado":
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		case "espaço já reservado para este horário":
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
		case "formato de data inválido. Use YYYY-MM-DD", 
			 "formato de hora inválido. Use HH:MM ou HH:MM:SS",
			 "a data não pode ser no passado",
			 "duração deve ser entre 1 e 24 horas",
			 "usuário não tem permissão para criar reservas",
			 "espaço não está disponível para reservas":
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar reserva", "details": errMsg})
		}
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
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
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservas [get]
func GetAllReservationsHandler(c *gin.Context) {
	reservations, err := services.GetAllReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar reservas"})
		return
	}
	c.JSON(http.StatusOK, reservations)
}
