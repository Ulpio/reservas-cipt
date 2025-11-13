package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Ulpio/reservas-cipt/dto"
	apierrors "github.com/Ulpio/reservas-cipt/errors"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// BuscarOuCriarClienteHandler busca um cliente pelo CPF ou cria um novo.
// @Summary Busca ou cria cliente
// @Description Procura um cliente pelo CPF e cria caso não exista.
// @Tags clientes
// @Accept json
// @Produce json
// @Param input body dto.ClienteInputDTO true "Dados do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/buscar-criar [post]
func BuscarOuCriarClienteHandler(c *gin.Context) {
	var input dto.ClienteInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		apierrors.SendError(c, apierrors.ErrDadosInvalidos.WithDetails("Verifique os campos obrigatórios"))
		return
	}

	cliente, err := services.BuscarOuCriarCliente(input)
	if err != nil {
		handleClientServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, cliente)
}

// BuscarClientePorCPF retorna os dados de um cliente a partir do CPF.
// @Summary Busca cliente por CPF
// @Description Retorna os dados de um cliente existente identificado pelo CPF.
// @Tags clientes
// @Produce json
// @Param cpf path string true "CPF do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /clientes/{cpf} [get]
func BuscarClientePorCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	cliente, err := services.GetClientByCPF(cpf)
	if err != nil {
		apierrors.SendError(c, apierrors.ErrClienteNaoEncontrado.WithDetails("Verifique se o CPF está correto"))
		return
	}
	c.JSON(http.StatusOK, cliente)
}

// GetAllClientes lista todos os clientes cadastrados.
// @Summary Lista clientes
// @Description Retorna todos os clientes registrados.
// @Tags clientes
// @Produce json
// @Success 200 {array} dto.ClienteOutputDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes [get]
func GetAllClientes(c *gin.Context) {
	clientes, err := services.GetAllClientes()
	if err != nil {
		apierrors.SendError(c, apierrors.ErrBancoDados.WithDetails("Erro ao listar clientes"))
		return
	}
	c.JSON(http.StatusOK, clientes)
}

// UpdateClientHandler atualiza os dados de um cliente.
// @Summary Atualiza cliente
// @Description Atualiza informações de um cliente existente.
// @Tags clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do cliente"
// @Param input body dto.ClienteInputDTO true "Novos dados do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/{id} [patch]
func UpdateClientHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apierrors.SendError(c, apierrors.ErrIDInvalido.WithDetails("O ID do cliente deve ser um número válido"))
		return
	}

	var input dto.ClienteInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		apierrors.SendError(c, apierrors.ErrDadosInvalidos.WithDetails("Verifique os campos fornecidos"))
		return
	}

	updated, err := services.UpdateClient(uint(id), input)
	if err != nil {
		handleClientServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// handleClientServiceError trata erros específicos dos serviços de cliente
func handleClientServiceError(c *gin.Context, err error) {
	errMsg := err.Error()

	// Erros de validação de CPF
	if strings.Contains(errMsg, "CPF inválido") || strings.Contains(errMsg, "CPF deve ter") {
		apierrors.SendError(c, apierrors.ErrCPFInvalido)
		return
	}
	if strings.Contains(errMsg, "CPF já cadastrado") {
		apierrors.SendError(c, apierrors.ErrCPFJaCadastrado)
		return
	}

	// Erros de validação de email
	if strings.Contains(errMsg, "email inválido") || strings.Contains(errMsg, "formato de email") {
		apierrors.SendError(c, apierrors.ErrEmailInvalido)
		return
	}
	if strings.Contains(errMsg, "email já cadastrado") {
		apierrors.SendError(c, apierrors.ErrEmailJaCadastrado)
		return
	}

	// Erros de validação de telefone
	if strings.Contains(errMsg, "telefone inválido") || strings.Contains(errMsg, "telefone deve conter") {
		apierrors.SendError(c, apierrors.ErrTelefoneInvalido)
		return
	}
	if strings.Contains(errMsg, "telefone já cadastrado") {
		apierrors.SendError(c, apierrors.ErrTelefoneJaCadastrado)
		return
	}

	// Erros de validação de idade
	if strings.Contains(errMsg, "18 anos") || strings.Contains(errMsg, "idade mínima") {
		apierrors.SendError(c, apierrors.ErrIdadeMinima)
		return
	}

	// Erros de validação de nome
	if strings.Contains(errMsg, "nome deve ter") || strings.Contains(errMsg, "nome inválido") {
		apierrors.SendError(c, apierrors.ErrNomeInvalido)
		return
	}

	// Erros de data de nascimento
	if strings.Contains(errMsg, "data de nascimento") {
		apierrors.SendError(c, apierrors.ErrDataNascimentoInvalida)
		return
	}

	// Erro genérico
	apierrors.SendError(c, apierrors.ErrInterno.WithDetails(errMsg))
}

