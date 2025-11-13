package services

import (
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
)

// GetDashboardStats retorna as estatísticas do dashboard
func GetDashboardStats() (dto.DashboardStatsDTO, error) {
	var stats dto.DashboardStatsDTO

	// Data de hoje (início e fim do dia)
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 1. Contar reservas de hoje
	var reservasHoje int64
	if err := database.DB.Model(&models.Reservation{}).
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Count(&reservasHoje).Error; err != nil {
		return stats, err
	}
	stats.ReservasHoje = reservasHoje

	// 2. Contar espaços ativos
	var espacosAtivos int64
	if err := database.DB.Model(&models.Space{}).
		Where("status = ?", "ativo").
		Count(&espacosAtivos).Error; err != nil {
		return stats, err
	}
	stats.EspacosAtivos = espacosAtivos

	// 3. Contar clientes ativos
	var clientesAtivos int64
	if err := database.DB.Model(&models.Client{}).
		Count(&clientesAtivos).Error; err != nil {
		return stats, err
	}
	stats.ClientesAtivos = clientesAtivos

	// 4. Calcular taxa de ocupação
	// Taxa = (Total de horas reservadas hoje / Total de horas disponíveis hoje) * 100
	// Horas disponíveis = Número de espaços ativos * 24 horas
	
	var totalHorasReservadas int64
	var result struct {
		Total int64
	}
	
	if err := database.DB.Model(&models.Reservation{}).
		Select("COALESCE(SUM(duration_hours), 0) as total").
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Scan(&result).Error; err != nil {
		return stats, err
	}
	totalHorasReservadas = result.Total

	// Calcular taxa de ocupação
	taxaOcupacao := 0
	if espacosAtivos > 0 {
		horasDisponiveis := espacosAtivos * 24
		if horasDisponiveis > 0 {
			taxaOcupacao = int((float64(totalHorasReservadas) / float64(horasDisponiveis)) * 100)
			// Garantir que não ultrapasse 100%
			if taxaOcupacao > 100 {
				taxaOcupacao = 100
			}
		}
	}
	stats.TaxaOcupacao = taxaOcupacao

	return stats, nil
}

