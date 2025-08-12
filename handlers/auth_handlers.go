package handlers

import (
	"net/http"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var input dto.LoginInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	token, err := services.AuthenticateUser(input.CPF, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erorr": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
