package handlers

import (
	"net/http"

	"github.com/Ulpio/reservas-cipt/services"
	"github.com/gin-gonic/gin"
)

// GetDashboardStatsHandler retorna as estatísticas do dashboard
// @Summary Obter estatísticas do dashboard
// @Description Retorna estatísticas principais para exibição no dashboard (reservas hoje, espaços ativos, clientes ativos, taxa de ocupação)
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.DashboardStatsDTO
// @Failure 401 {object} dto.ErrorResponse "Não autenticado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /dashboard/stats [get]
func GetDashboardStatsHandler(c *gin.Context) {
	stats, err := services.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar estatísticas"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

