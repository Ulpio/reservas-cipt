package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	cliente, err := services.BuscarOuCriarCliente(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ou criar cliente"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

// BuscarClientePorCPF retorna os dados de um cliente a partir do CPF.
// @Summary Busca cliente por CPF
// @Description Retorna os dados de um cliente existente identificado pelo CPF.
// @Tags clientes
// @Security BearerAuth
// @Produce json
// @Param cpf path string true "CPF do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /clientes/cpf/{cpf} [get]
func BuscarClientePorCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	cliente, err := services.GetClientByCPF(cpf)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}
	c.JSON(http.StatusOK, cliente)
}

// BuscarClientePorTelefone retorna os dados de um cliente a partir do telefone.
// @Summary Busca cliente por telefone
// @Description Retorna os dados de um cliente existente identificado pelo telefone.
// @Tags clientes
// @Security BearerAuth
// @Produce json
// @Param phone path string true "Telefone do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /clientes/phone/{phone} [get]
func BuscarClientePorTelefone(c *gin.Context) {
	phone := c.Param("phone")
	cliente, err := services.GetClientByPhone(phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}
	c.JSON(http.StatusOK, cliente)
}

// BuscarClientePorID retorna os dados de um cliente a partir do ID.
// @Summary Busca cliente por ID
// @Description Retorna os dados de um cliente existente identificado pelo ID.
// @Tags clientes
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /clientes/{id} [get]
func BuscarClientePorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	cliente, err := services.GetClientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}
	c.JSON(http.StatusOK, cliente)
}

// CreateClientHandler cria um novo cliente.
// @Summary Cria um novo cliente
// @Description Cadastra um novo cliente no sistema com validações completas.
// @Tags clientes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.ClienteInputDTO true "Dados do cliente"
// @Success 201 {object} dto.ClienteOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes [post]
func CreateClientHandler(c *gin.Context) {
	var input dto.ClienteInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	cliente, err := services.CreateClient(input)
	if err != nil {
		errMsg := err.Error()
		
		// Retornar status code apropriado baseado no erro
		switch errMsg {
		case "CPF já cadastrado no sistema", "Telefone já cadastrado":
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
		case "CPF deve ter 11 dígitos", 
			 "CPF inválido",
			 "formato de email inválido",
			 "telefone deve ter entre 10 e 11 dígitos",
			 "formato de data inválido. Use YYYY-MM-DD",
			 "data de nascimento deve ser uma data passada",
			 "cliente deve ter pelo menos 18 anos":
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar cliente", "details": errMsg})
		}
		return
	}

	c.JSON(http.StatusCreated, cliente)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ou criar cliente"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input dto.ClienteInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := services.UpdateClient(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar cliente"})
		return
	}

	c.JSON(http.StatusOK, updated)
}
