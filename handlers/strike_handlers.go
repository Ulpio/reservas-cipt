package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// CreateStrikeHandler registra um novo strike para um cliente.
// @Summary Registra strike
// @Description Cria uma advertência para um cliente.
// @Tags strikes
// @Accept json
// @Produce json
// @Param input body dto.StrikeInputDTO true "Dados do strike"
// @Success 201 {object} dto.StrikeOutputDTO
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /strikes [post]
func CreateStrikeHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "recepcionista" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para adicionar strikes"})
		return
	}
	var input dto.StrikeInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input inválido"})
		return
	}
	strike, err := services.CreateStrike(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, strike)
}

// GetStrikesByClientHandler lista os strikes de um cliente.
// @Summary Lista strikes por cliente
// @Description Retorna as advertências associadas a um cliente específico.
// @Tags strikes
// @Produce json
// @Param id path int true "ID do cliente"
// @Success 200 {array} dto.StrikeOutputDTO
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /strikes/client/{id} [get]
func GetStrikesByClientHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" && role != "recepcionista" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para listar os strikes"})
		return
	}
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	strikes, err := services.GetStrikesByClient(uint(clientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, strikes)
}

// RevokeStrikeHandler revoga um strike existente.
// @Summary Revoga strike
// @Description Revoga uma advertência de um cliente.
// @Tags strikes
// @Param id path int true "ID do strike"
// @Success 200 {object} dto.StrikeOutputDTO
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /strikes/{id} [delete]
func RevokeStrikeHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para revogar strikes"})
		return
	}
	id := c.Param("id")
	strikeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	strike, err := services.RevokeStrike(uint(strikeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, strike)
}
