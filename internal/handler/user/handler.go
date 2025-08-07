package handler

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/internal/dto"
	"github.com/Ulpio/reservas-cipt/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	UseCase *user.UserUseCase
}

func NewHandler(uc *user.UserUseCase) *UserHandler {
	return &UserHandler{UseCase: uc}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.UseCase.RegisterUser(input.Name, input.CPF, input.Role, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "usuario criado com sucesso"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	user, err := h.UseCase.FindUserByID(uint(userID))
	if err != nil && err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurso não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UseCase.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) DeleteByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}
	err = h.UseCase.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusNoContent)
}
