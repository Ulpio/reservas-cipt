package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/routes"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/Ulpio/reservas-cipt/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupAuthRoutes(r)
	routes.SetupUserRoutes(r)
	return r
}

func TestCreateUserHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	router := setupTestRouter()

	body := dto.UserInputDTO{
		Name: "Test User", CPF: "12345678900", Role: "admin",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetUserByIDHandler(t *testing.T) {
	tests.SetupTestDB(t)

	user, _ := services.CreateUser("Me User", "22233344455", "admin")
	token, _ := utils.GenerateJWT(user.ID, user.Role)

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(int(user.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Me User")
}

func TestDeleteUserHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "11122233344", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	user, _ := services.CreateUser("To Delete", "99988877766", "admin")

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodDelete, "/users/"+strconv.Itoa(int(user.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestMeHandler(t *testing.T) {
	tests.SetupTestDB(t)

	// cria usuário e obtém token
	user, _ := services.CreateUser("Me User", "22233344455", "admin")
	token, _ := utils.GenerateJWT(user.ID, user.Role)

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Me User")
}
