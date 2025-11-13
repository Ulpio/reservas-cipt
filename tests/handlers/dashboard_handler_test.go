package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/Ulpio/reservas-cipt/utils"
	"github.com/gin-gonic/gin"
)

func TestGetDashboardStats(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM users")
	database.DB.Exec("DELETE FROM spaces")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM reservations")
	database.DB.Exec("DELETE FROM strikes")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.JWTAuthMiddleware())
	{
		dashboard.GET("/stats", handlers.GetDashboardStatsHandler)
	}

	// Criar usuário admin
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	// Criar dados de teste
	// 1. Criar espaços
	space1 := models.Space{Name: "Sala 1", Type: "visitantes", Status: "ativo", Capacity: 10}
	space2 := models.Space{Name: "Sala 2", Type: "visitantes", Status: "ativo", Capacity: 20}
	space3 := models.Space{Name: "Sala 3", Type: "visitantes", Status: "manutencao", Capacity: 15}
	database.DB.Create(&space1)
	database.DB.Create(&space2)
	database.DB.Create(&space3)

	// 2. Criar clientes
	client1 := models.Client{Name: "Cliente 1", CPF: "11111111111", Email: "c1@test.com", Phone: "1111111111"}
	client2 := models.Client{Name: "Cliente 2", CPF: "22222222222", Email: "c2@test.com", Phone: "2222222222"}
	client3 := models.Client{Name: "Cliente 3", CPF: "33333333333", Email: "c3@test.com", Phone: "3333333333"}
	database.DB.Create(&client1)
	database.DB.Create(&client2)
	database.DB.Create(&client3)

	// 3. Criar reservas de hoje
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tomorrow := today.AddDate(0, 0, 1)
	startTime1 := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local)
	startTime2 := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.Local)
	startTime3 := time.Date(now.Year(), now.Month(), now.Day()+1, 10, 0, 0, 0, time.Local)

	reservation1 := models.Reservation{
		ClientID:        client1.ID,
		ReceptionistID:  admin.ID,
		SpaceID:         space1.ID,
		Date:            today,
		StartTime:       startTime1,
		DurationHours:   2,
	}
	reservation2 := models.Reservation{
		ClientID:        client2.ID,
		ReceptionistID:  admin.ID,
		SpaceID:         space2.ID,
		Date:            today,
		StartTime:       startTime2,
		DurationHours:   3,
	}
	reservation3 := models.Reservation{
		ClientID:        client3.ID,
		ReceptionistID:  admin.ID,
		SpaceID:         space1.ID,
		Date:            tomorrow,
		StartTime:       startTime3,
		DurationHours:   4,
	}
	database.DB.Create(&reservation1)
	database.DB.Create(&reservation2)
	database.DB.Create(&reservation3)

	// Testar endpoint
	req, _ := http.NewRequest("GET", "/dashboard/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperado status 200, recebido %d", w.Code)
	}

	var stats dto.DashboardStatsDTO
	if err := json.Unmarshal(w.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	// Verificar valores esperados
	if stats.ReservasHoje != 2 {
		t.Errorf("Esperado 2 reservas hoje, recebido %d", stats.ReservasHoje)
	}

	if stats.EspacosAtivos != 2 {
		t.Errorf("Esperado 2 espaços ativos, recebido %d", stats.EspacosAtivos)
	}

	if stats.ClientesAtivos != 3 {
		t.Errorf("Esperado 3 clientes ativos, recebido %d", stats.ClientesAtivos)
	}

	// Taxa de ocupação = (2+3) / (2*24) * 100 = 5 / 48 * 100 ≈ 10%
	if stats.TaxaOcupacao < 0 || stats.TaxaOcupacao > 100 {
		t.Errorf("Taxa de ocupação fora do intervalo válido: %d", stats.TaxaOcupacao)
	}
}

func TestGetDashboardStatsUnauthorized(t *testing.T) {
	tests.SetupTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.JWTAuthMiddleware())
	{
		dashboard.GET("/stats", handlers.GetDashboardStatsHandler)
	}

	req, _ := http.NewRequest("GET", "/dashboard/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Esperado status 401, recebido %d", w.Code)
	}
}

