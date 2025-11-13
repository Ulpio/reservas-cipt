package dto

// DashboardStatsDTO representa as estatísticas do dashboard
type DashboardStatsDTO struct {
	ReservasHoje    int64 `json:"reservasHoje"`
	EspacosAtivos   int64 `json:"espacosAtivos"`
	ClientesAtivos  int64 `json:"clientesAtivos"`
	TaxaOcupacao    int   `json:"taxaOcupacao"`
}

