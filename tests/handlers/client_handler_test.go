package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestCreateClient(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	clientes := r.Group("/clientes")
	clientes.Use(middleware.JWTAuthMiddleware())
	{
		clientes.POST("", handlers.CreateClientHandler)
	}

	// Criar admin e token
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	// Teste 1: Criar cliente com sucesso
	input := dto.ClienteInputDTO{
		Name:  "João Silva",
		CPF:   "123.456.789-00",
		Phone: "(11) 98765-4321",
		Email: "joao@email.com",
	}
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest("POST", "/clientes", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Esperado status 201, recebido %d: %s", w.Code, w.Body.String())
	}

	var response dto.ClienteOutputDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "João Silva" {
		t.Errorf("Esperado nome 'João Silva', recebido '%s'", response.Name)
	}

	// CPF deve ser normalizado (sem pontos e traços)
	if response.CPF != "12345678900" {
		t.Errorf("Esperado CPF normalizado '12345678900', recebido '%s'", response.CPF)
	}

	// Teste 2: Tentar criar cliente com CPF duplicado
	req2, _ := http.NewRequest("POST", "/clientes", bytes.NewBuffer(body))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Esperado status 400, recebido %d", w2.Code)
	}

	var errorResp map[string]string
	json.Unmarshal(w2.Body.Bytes(), &errorResp)
	if errorResp["error"] != "CPF já cadastrado" {
		t.Errorf("Esperado erro 'CPF já cadastrado', recebido '%s'", errorResp["error"])
	}
}

func TestGetClientByCPF(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	clientes := r.Group("/clientes")
	clientes.Use(middleware.JWTAuthMiddleware())
	{
		clientes.GET("/cpf/:cpf", handlers.BuscarClientePorCPF)
	}

	// Criar admin e token
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	// Criar cliente
	client := models.Client{
		Name:  "Maria Santos",
		CPF:   "98765432100",
		Phone: "11987654321",
		Email: "maria@email.com",
	}
	database.DB.Create(&client)

	// Teste 1: Buscar com CPF formatado
	req, _ := http.NewRequest("GET", "/clientes/cpf/987.654.321-00", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperado status 200, recebido %d: %s", w.Code, w.Body.String())
	}

	var response dto.ClienteOutputDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "Maria Santos" {
		t.Errorf("Esperado nome 'Maria Santos', recebido '%s'", response.Name)
	}

	// Teste 2: Buscar CPF não existente
	req2, _ := http.NewRequest("GET", "/clientes/cpf/11111111111", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Esperado status 404, recebido %d", w2.Code)
	}
}

func TestGetClientByPhone(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	clientes := r.Group("/clientes")
	clientes.Use(middleware.JWTAuthMiddleware())
	{
		clientes.GET("/phone/:phone", handlers.BuscarClientePorTelefone)
	}

	// Criar admin e token
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	// Criar cliente
	client := models.Client{
		Name:  "Pedro Costa",
		CPF:   "11122233344",
		Phone: "21987654321",
		Email: "pedro@email.com",
	}
	database.DB.Create(&client)

	// Teste: Buscar com telefone formatado
	req, _ := http.NewRequest("GET", "/clientes/phone/(21) 98765-4321", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperado status 200, recebido %d: %s", w.Code, w.Body.String())
	}

	var response dto.ClienteOutputDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "Pedro Costa" {
		t.Errorf("Esperado nome 'Pedro Costa', recebido '%s'", response.Name)
	}
}

func TestGetClientByID(t *testing.T) {
	tests.SetupTestDB(t)
	database.DB.Exec("DELETE FROM clients")
	database.DB.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	clientes := r.Group("/clientes")
	clientes.Use(middleware.JWTAuthMiddleware())
	{
		clientes.GET("/:id", handlers.BuscarClientePorID)
	}

	// Criar admin e token
	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	// Criar cliente
	client := models.Client{
		Name:  "Ana Paula",
		CPF:   "55566677788",
		Phone: "31987654321",
		Email: "ana@email.com",
	}
	database.DB.Create(&client)

	// Teste: Buscar por ID
	req, _ := http.NewRequest("GET", "/clientes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperado status 200, recebido %d: %s", w.Code, w.Body.String())
	}

	var response dto.ClienteOutputDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "Ana Paula" {
		t.Errorf("Esperado nome 'Ana Paula', recebido '%s'", response.Name)
	}
}

func TestClientUnauthorized(t *testing.T) {
	tests.SetupTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	clientes := r.Group("/clientes")
	clientes.Use(middleware.JWTAuthMiddleware())
	{
		clientes.GET("/cpf/:cpf", handlers.BuscarClientePorCPF)
	}

	// Teste sem token
	req, _ := http.NewRequest("GET", "/clientes/cpf/12345678900", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Esperado status 401, recebido %d", w.Code)
	}
}

