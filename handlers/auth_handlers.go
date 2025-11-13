package handlers

import (
	"net/http"

	"github.com/Ulpio/reservas-cipt/dto"
	apierrors "github.com/Ulpio/reservas-cipt/errors"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// LoginHandler autentica um usuário e retorna um token JWT.
// @Summary Autentica um usuário
// @Description Endpoint público utilizado para autenticar usuários do sistema.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginInputDTO true "Credenciais de acesso"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var input dto.LoginInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		apierrors.SendError(c, apierrors.ErrDadosInvalidos.WithDetails("Verifique se CPF e senha foram fornecidos"))
		return
	}
	token, err := services.AuthenticateUser(input.CPF, input.Password)
	if err != nil {
		apierrors.SendError(c, apierrors.ErrCredenciaisInvalidas)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
