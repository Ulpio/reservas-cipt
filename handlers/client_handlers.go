package handlers

import (
	"net/http"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

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

func BuscarClientePorCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	cliente, err := services.GetClientByCPF(cpf)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario nao encontrado", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cliente)
}

func GetAllClientes(c *gin.Context) {
	clientes, err := services.GetAllClientes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ou criar cliente"})
		return
	}
	c.JSON(http.StatusOK, clientes)
}
