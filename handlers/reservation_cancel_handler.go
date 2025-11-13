package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// CancelReservationHandler cancela (deleta) uma reserva
// @Summary Cancela uma reserva
// @Description Remove uma reserva e envia email de cancelamento para o cliente
// @Tags reservas
// @Security BearerAuth
// @Param id path int true "ID da reserva"
// @Success 200 {object} map[string]string "message: Reserva cancelada com sucesso"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reservas/{id} [delete]
func CancelReservationHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se a reserva existe antes de cancelar
	_, err = services.GetReservationDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva não encontrada"})
		return
	}

	// Cancelar a reserva
	if err := services.CancelReservation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cancelar reserva"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reserva cancelada com sucesso"})
}

