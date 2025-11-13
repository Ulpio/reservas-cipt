package handlers_test

import (
	"bytes"
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

func TestCreateReservationWithStringFormats(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM reservations")
	database.DB.Exec("DELETE FROM spaces")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	reservas := r.Group("/reservas")
	reservas.Use(middleware.JWTAuthMiddleware())
	{
		reservas.POST("", handlers.CreateReservationHandler)
	}

	// Criar dados de teste
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	client := models.Client{Name: "João Silva", CPF: "12345678900", Phone: "11987654321", Email: "joao@test.com"}
	database.DB.Create(&client)

	space := models.Space{Name: "Sala A", Type: "visitantes", Status: "ativo", Capacity: 10}
	database.DB.Create(&space)

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Teste 1: Criar reserva com HH:MM
	input1 := dto.CreateReservationInputDTO{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           tomorrow,
		StartTime:      "16:00",
		DurationHours:  2,
	}
	body1, _ := json.Marshal(input1)

	req1, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(body1))
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("Esperado status 201, recebido %d: %s", w1.Code, w1.Body.String())
	}

	// Teste 2: Criar reserva com HH:MM:SS
	input2 := dto.CreateReservationInputDTO{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           tomorrow,
		StartTime:      "19:00:00",
		DurationHours:  1,
	}
	body2, _ := json.Marshal(input2)

	req2, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(body2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusCreated {
		t.Errorf("Esperado status 201, recebido %d: %s", w2.Code, w2.Body.String())
	}
}

func TestCreateReservationValidations(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM reservations")
	database.DB.Exec("DELETE FROM spaces")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	reservas := r.Group("/reservas")
	reservas.Use(middleware.JWTAuthMiddleware())
	{
		reservas.POST("", handlers.CreateReservationHandler)
	}

	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	client := models.Client{Name: "Maria Santos", CPF: "98765432100", Phone: "21987654321", Email: "maria@test.com"}
	database.DB.Create(&client)

	space := models.Space{Name: "Sala B", Type: "visitantes", Status: "ativo", Capacity: 20}
	database.DB.Create(&space)

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Teste 1: Data no passado
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	inputPast := dto.CreateReservationInputDTO{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           yesterday,
		StartTime:      "10:00",
		DurationHours:  2,
	}
	bodyPast, _ := json.Marshal(inputPast)

	reqPast, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(bodyPast))
	reqPast.Header.Set("Authorization", "Bearer "+token)
	reqPast.Header.Set("Content-Type", "application/json")
	wPast := httptest.NewRecorder()
	r.ServeHTTP(wPast, reqPast)

	if wPast.Code != http.StatusBadRequest {
		t.Errorf("Data passada deveria retornar 400, recebido %d", wPast.Code)
	}

	// Teste 2: Cliente não existente
	inputNoClient := dto.CreateReservationInputDTO{
		ClientID:       999999,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           tomorrow,
		StartTime:      "10:00",
		DurationHours:  2,
	}
	bodyNoClient, _ := json.Marshal(inputNoClient)

	reqNoClient, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(bodyNoClient))
	reqNoClient.Header.Set("Authorization", "Bearer "+token)
	reqNoClient.Header.Set("Content-Type", "application/json")
	wNoClient := httptest.NewRecorder()
	r.ServeHTTP(wNoClient, reqNoClient)

	if wNoClient.Code != http.StatusNotFound {
		t.Errorf("Cliente não existente deveria retornar 404, recebido %d", wNoClient.Code)
	}

	// Teste 3: Espaço em manutenção
	spaceManutencao := models.Space{Name: "Sala C", Type: "visitantes", Status: "manutencao", Capacity: 15}
	database.DB.Create(&spaceManutencao)

	inputManutencao := dto.CreateReservationInputDTO{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        spaceManutencao.ID,
		Date:           tomorrow,
		StartTime:      "10:00",
		DurationHours:  2,
	}
	bodyManutencao, _ := json.Marshal(inputManutencao)

	reqManutencao, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(bodyManutencao))
	reqManutencao.Header.Set("Authorization", "Bearer "+token)
	reqManutencao.Header.Set("Content-Type", "application/json")
	wManutencao := httptest.NewRecorder()
	r.ServeHTTP(wManutencao, reqManutencao)

	if wManutencao.Code != http.StatusBadRequest {
		t.Errorf("Espaço em manutenção deveria retornar 400, recebido %d", wManutencao.Code)
	}
}

func TestCreateReservationConflict(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM reservations")
	database.DB.Exec("DELETE FROM spaces")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	reservas := r.Group("/reservas")
	reservas.Use(middleware.JWTAuthMiddleware())
	{
		reservas.POST("", handlers.CreateReservationHandler)
	}

	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	client := models.Client{Name: "Pedro Costa", CPF: "11122233344", Phone: "31987654321", Email: "pedro@test.com"}
	database.DB.Create(&client)

	space := models.Space{Name: "Sala D", Type: "visitantes", Status: "ativo", Capacity: 10}
	database.DB.Create(&space)

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	tomorrowTime, _ := time.ParseInLocation("2006-01-02", tomorrow, time.Local)
	startTime := time.Date(tomorrowTime.Year(), tomorrowTime.Month(), tomorrowTime.Day(), 14, 0, 0, 0, time.Local)

	// Criar reserva existente
	existingReservation := models.Reservation{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           tomorrowTime,
		StartTime:      startTime,
		DurationHours:  2, // 14:00 - 16:00
	}
	database.DB.Create(&existingReservation)

	// Tentar criar reserva conflitante
	inputConflict := dto.CreateReservationInputDTO{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           tomorrow,
		StartTime:      "15:00", // Conflita com 14:00-16:00
		DurationHours:  2,
	}
	bodyConflict, _ := json.Marshal(inputConflict)

	reqConflict, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(bodyConflict))
	reqConflict.Header.Set("Authorization", "Bearer "+token)
	reqConflict.Header.Set("Content-Type", "application/json")
	wConflict := httptest.NewRecorder()
	r.ServeHTTP(wConflict, reqConflict)

	if wConflict.Code != http.StatusConflict {
		t.Errorf("Reserva conflitante deveria retornar 409, recebido %d: %s", wConflict.Code, wConflict.Body.String())
	}
}

func TestCreateReservationSameDay(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM reservations")
	database.DB.Exec("DELETE FROM spaces")
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	reservas := r.Group("/reservas")
	reservas.Use(middleware.JWTAuthMiddleware())
	{
		reservas.POST("", handlers.CreateReservationHandler)
	}

	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	client := models.Client{Name: "Ana Paula", CPF: "55566677788", Phone: "41987654321", Email: "ana@test.com"}
	database.DB.Create(&client)

	space := models.Space{Name: "Sala E", Type: "visitantes", Status: "ativo", Capacity: 15}
	database.DB.Create(&space)

	// Usar a data de hoje
	today := time.Now().Format("2006-01-02")

	// Criar reserva para hoje (simula: são 15:30 e a reserva é para 16:00)
	inputToday := dto.CreateReservationInputDTO{
		ClientID:       client.ID,
		ReceptionistID: admin.ID,
		SpaceID:        space.ID,
		Date:           today,
		StartTime:      "16:00",
		DurationHours:  1,
	}
	bodyToday, _ := json.Marshal(inputToday)

	reqToday, _ := http.NewRequest("POST", "/reservas", bytes.NewBuffer(bodyToday))
	reqToday.Header.Set("Authorization", "Bearer "+token)
	reqToday.Header.Set("Content-Type", "application/json")
	wToday := httptest.NewRecorder()
	r.ServeHTTP(wToday, reqToday)

	if wToday.Code != http.StatusCreated {
		t.Errorf("Reserva para hoje deveria retornar 201, recebido %d: %s", wToday.Code, wToday.Body.String())
	}

	var response dto.ReservationOutputDTO
	json.Unmarshal(wToday.Body.Bytes(), &response)

	if response.ID == 0 {
		t.Error("Reserva não foi criada com sucesso")
	}
}

