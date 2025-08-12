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
	router := setupTestRouter()

	body := dto.UserInputDTO{
		Name: "Test User", CPF: "12345678900", Role: "admin",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetUserByIDHandler(t *testing.T) {
	user, _ := services.CreateUser("Me User", "22233344455", "admin")
	token, _ := utils.GenerateJWT(user.ID, user.Role)

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(int(user.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Handler Test")
}

func TestDeleteUserHandler(t *testing.T) {
	user, _ := services.CreateUser("To Delete", "99988877766", "admin")

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodDelete, "/users/"+strconv.Itoa(int(user.ID)), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestMeHandler(t *testing.T) {
	// cria usuário e obtém token
	user, _ := services.CreateUser("Me User", "22233344455", "admin")
	token, _ := utils.GenerateJWT(user.ID, user.Role)

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Me User")
}
