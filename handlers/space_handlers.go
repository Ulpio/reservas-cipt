package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// CreateSpaceHandler cria um novo espaço disponível para reserva.
// @Summary Cria um novo espaço físico
// @Description Endpoint utilizado por administradores para cadastrar um novo espaço disponível para reserva.
// @Tags espacos
// @Accept json
// @Produce json
// @Param input body dto.CreateSpaceDTO true "Dados do espaço a ser criado"
// @Success 201 {object} dto.SpaceOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /espacos [post]
func CreateSpaceHandler(c *gin.Context) {
	var input dto.CreateSpaceDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	space, err := services.CreateSpace(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, space)
}

// GetAllSpacesHandler lista todos os espaços cadastrados.
// @Summary Lista espaços
// @Description Retorna todos os espaços registrados.
// @Tags espacos
// @Produce json
// @Success 200 {array} dto.SpaceOutputDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos [get]
func GetAllSpacesHandler(c *gin.Context) {
	spaces, err := services.GetAllSpaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, spaces)
}

// GetSpacesByIDHandler retorna um espaço específico pelo ID.
// @Summary Busca espaço por ID
// @Description Obtém os dados de um espaço a partir do seu identificador.
// @Tags espacos
// @Produce json
// @Param id path int true "ID do espaço"
// @Success 200 {object} dto.SpaceOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos/{id} [get]
func GetSpacesByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	sala, err := services.GetSpaceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if sala.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "sala não encontrada"})
		return
	}
	c.JSON(http.StatusOK, sala)
}

// UpdateSpaceHandler atualiza todos os campos de um espaço.
// @Summary Atualiza espaço
// @Description Atualiza as informações completas de um espaço existente.
// @Tags espacos
// @Accept json
// @Produce json
// @Param id path int true "ID do espaço"
// @Param input body dto.UpdateSpaceDTO true "Novos dados do espaço"
// @Success 200 {object} dto.SpaceOutputDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos/{id} [put]
func UpdateSpaceHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input dto.UpdateSpaceDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := services.UpdateSpace(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar espaço"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSpaceHandler exclui um espaço do sistema.
// @Summary Exclui espaço
// @Description Remove um espaço previamente cadastrado.
// @Tags espacos
// @Param id path int true "ID do espaço"
// @Success 204 {object} nil
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos/{id} [delete]
func DeleteSpaceHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := services.DeleteSpace(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar espaço"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateSpaceStatusHandler atualiza apenas o status de um espaço.
// @Summary Atualiza status do espaço
// @Description Altera o status operacional de um espaço.
// @Tags espacos
// @Accept json
// @Produce json
// @Param id path int true "ID do espaço"
// @Param input body dto.UpdateStatusDTO true "Novo status"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos/{id}/status [patch]
func UpdateSpaceStatusHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input dto.UpdateStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateSpaceStatus(uint(id), input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Status atualizado com sucesso"})
}

// UpdateSpaceNoticeHandler atualiza apenas o aviso de um espaço.
// @Summary Atualiza aviso do espaço
// @Description Modifica o aviso exibido para um espaço.
// @Tags espacos
// @Accept json
// @Produce json
// @Param id path int true "ID do espaço"
// @Param input body dto.UpdateNoticeDTO true "Novo aviso"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /espacos/{id}/aviso [patch]
func UpdateSpaceNoticeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input dto.UpdateNoticeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateSpaceNotice(uint(id), input.Notice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar aviso"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Aviso atualizado com sucesso"})
}
