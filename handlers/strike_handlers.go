package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

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
