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
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
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
// @Produce json
// @Param cpf path string true "CPF do cliente"
// @Success 200 {object} dto.ClienteOutputDTO
// @Failure 404 {object} gin.H
// @Router /clientes/{cpf} [get]
func BuscarClientePorCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	cliente, err := services.GetClientByCPF(cpf)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario nao encontrado", "details": err.Error()})
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
// @Failure 500 {object} gin.H
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
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
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
