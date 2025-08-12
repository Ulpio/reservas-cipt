package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	// cria um usuário para login
	createdUser, _ := services.CreateUser("Login Test", "00011122233", "admin")

	router := setupTestRouter()
	login := dto.LoginInputDTO{CPF: "00011122233", Password: "12345678"}
	jsonLogin, _ := json.Marshal(login)

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("Criando usuário:", createdUser)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}
