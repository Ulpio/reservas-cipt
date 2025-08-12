package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

func GetAllUsersHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para listar os usuários"})
		return
	}
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByIDHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para listar os usuários"})
		return
	}
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id invalido"})
		return
	}
	users, err := services.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUserHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para registrar"})
		return
	}
	var input dto.UserInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}
	newUser, err := services.CreateUser(input.Name, input.CPF, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

func DeleteUserHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "você não tem autorização para registrar"})
		return
	}
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id invalido"})
		return
	}
	err = services.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func MeHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de usuário inválido"})
		return
	}

	userDTO, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, userDTO)
}
