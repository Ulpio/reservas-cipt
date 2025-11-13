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

func TestGetAllSpacesStatus(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM reservations")
	database.DB.Exec("DELETE FROM spaces")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	espacos := r.Group("/espacos")
	espacos.Use(middleware.JWTAuthMiddleware())
	{
		espacos.GET("/status", handlers.GetAllSpacesStatusHandler)
	}

	// Criar admin e token
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	// Criar espaços
	space1 := models.Space{Name: "Sala Disponível", Type: "visitantes", Status: "ativo", Capacity: 10}
	space2 := models.Space{Name: "Sala Ocupada", Type: "visitantes", Status: "ativo", Capacity: 20}
	space3 := models.Space{Name: "Sala Manutenção", Type: "visitantes", Status: "manutencao", Capacity: 15}
	database.DB.Create(&space1)
	database.DB.Create(&space2)
	database.DB.Create(&space3)

	// Criar cliente
	client := models.Client{Name: "João Silva", CPF: "12345678900", Phone: "11987654321", Email: "joao@test.com"}
	database.DB.Create(&client)

	// Criar reserva atual (em andamento)
	now := time.Now()
	startTime := now.Add(-1 * time.Hour) // Começou há 1 hora
	dateToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	reservationCurrent := models.Reservation{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space2.ID,
		Date:           dateToday,
		StartTime:      startTime,
		DurationHours:  3, // Vai até daqui 2 horas
	}
	database.DB.Create(&reservationCurrent)

	// Criar reserva futura (próxima)
	startTimeFuture := now.Add(2 * time.Hour)
	reservationFuture := models.Reservation{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space2.ID,
		Date:           dateToday,
		StartTime:      startTimeFuture,
		DurationHours:  2,
	}
	database.DB.Create(&reservationFuture)

	// Criar reserva no espaço disponível (no futuro)
	startTimeFuture2 := now.Add(5 * time.Hour)
	reservationFuture2 := models.Reservation{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space1.ID,
		Date:           dateToday,
		StartTime:      startTimeFuture2,
		DurationHours:  1,
	}
	database.DB.Create(&reservationFuture2)

	// Testar endpoint
	req, _ := http.NewRequest("GET", "/espacos/status", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperado status 200, recebido %d: %s", w.Code, w.Body.String())
	}

	var response []dto.SpaceStatusDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Esperado 3 espaços, recebido %d", len(response))
	}

	// Verificar espaço disponível
	var salaDisponivel *dto.SpaceStatusDTO
	for i := range response {
		if response[i].Nome == "Sala Disponível" {
			salaDisponivel = &response[i]
			break
		}
	}

	if salaDisponivel == nil {
		t.Fatal("Sala Disponível não encontrada na resposta")
	}

	if salaDisponivel.Status != "disponivel" {
		t.Errorf("Esperado status 'disponivel', recebido '%s'", salaDisponivel.Status)
	}

	if salaDisponivel.ReservaAtual != nil {
		t.Errorf("Sala disponível não deveria ter reserva atual")
	}

	if salaDisponivel.ProximaReserva == nil {
		t.Errorf("Sala disponível deveria ter próxima reserva")
	}

	if len(salaDisponivel.ReservasHoje) != 1 {
		t.Errorf("Esperado 1 reserva hoje, recebido %d", len(salaDisponivel.ReservasHoje))
	}

	// Verificar espaço ocupado
	var salaOcupada *dto.SpaceStatusDTO
	for i := range response {
		if response[i].Nome == "Sala Ocupada" {
			salaOcupada = &response[i]
			break
		}
	}

	if salaOcupada == nil {
		t.Fatal("Sala Ocupada não encontrada na resposta")
	}

	if salaOcupada.Status != "ocupado" {
		t.Errorf("Esperado status 'ocupado', recebido '%s'", salaOcupada.Status)
	}

	if salaOcupada.ReservaAtual == nil {
		t.Errorf("Sala ocupada deveria ter reserva atual")
	} else {
		if salaOcupada.ReservaAtual.Cliente != "João Silva" {
			t.Errorf("Esperado cliente 'João Silva', recebido '%s'", salaOcupada.ReservaAtual.Cliente)
		}
	}

	if salaOcupada.ProximaReserva == nil {
		t.Errorf("Sala ocupada deveria ter próxima reserva")
	}

	if len(salaOcupada.ReservasHoje) != 2 {
		t.Errorf("Esperado 2 reservas hoje, recebido %d", len(salaOcupada.ReservasHoje))
	}

	// Verificar espaço em manutenção
	var salaManutencao *dto.SpaceStatusDTO
	for i := range response {
		if response[i].Nome == "Sala Manutenção" {
			salaManutencao = &response[i]
			break
		}
	}

	if salaManutencao == nil {
		t.Fatal("Sala Manutenção não encontrada na resposta")
	}

	if salaManutencao.Status != "manutencao" {
		t.Errorf("Esperado status 'manutencao', recebido '%s'", salaManutencao.Status)
	}

	if salaManutencao.ReservaAtual != nil {
		t.Errorf("Sala em manutenção não deveria ter reserva atual")
	}

	if len(salaManutencao.ReservasHoje) != 0 {
		t.Errorf("Esperado 0 reservas hoje, recebido %d", len(salaManutencao.ReservasHoje))
	}
}

func TestGetAllSpacesStatusUnauthorized(t *testing.T) {
	tests.SetupTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	espacos := r.Group("/espacos")
	espacos.Use(middleware.JWTAuthMiddleware())
	{
		espacos.GET("/status", handlers.GetAllSpacesStatusHandler)
	}

	req, _ := http.NewRequest("GET", "/espacos/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Esperado status 401, recebido %d", w.Code)
	}
}

