package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// GetAllUsersHandler lista todos os usuários cadastrados.
// @Summary Lista usuários
// @Description Endpoint restrito a administradores para listagem de usuários.
// @Tags usuarios
// @Produce json
// @Success 200 {array} dto.UserOutputDTO
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [get]
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

// GetUserByIDHandler retorna um usuário específico pelo ID.
// @Summary Busca usuário por ID
// @Description Recupera os dados de um usuário específico.
// @Tags usuarios
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} dto.UserOutputDTO
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /users/{id} [get]
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

// CreateUserHandler cria um novo usuário administrador ou recepcionista.
// @Summary Cria usuário
// @Description Endpoint restrito a administradores para cadastro de novos usuários.
// @Tags usuarios
// @Accept json
// @Produce json
// @Param input body dto.UserInputDTO true "Dados do usuário"
// @Success 201 {object} dto.UserOutputDTO
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
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

// DeleteUserHandler remove um usuário existente.
// @Summary Remove usuário
// @Description Exclui um usuário existente.
// @Tags usuarios
// @Param id path int true "ID do usuário"
// @Success 204 {object} nil
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/{id} [delete]
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

// MeHandler retorna os dados do usuário autenticado.
// @Summary Retorna usuário autenticado
// @Description Recupera informações do usuário baseado no token de autenticação.
// @Tags usuarios
// @Produce json
// @Success 200 {object} dto.UserOutputDTO
// @Failure 401 {object} gin.H
// @Router /users/me [get]
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
