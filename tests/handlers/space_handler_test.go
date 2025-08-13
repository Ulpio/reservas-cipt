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

func setupSpaceRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SpaceRoutes(r)
	return r
}

func TestCreateSpaceHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	router := setupSpaceRouter()

	body := dto.CreateSpaceDTO{Name: "Sala 1", Type: "sala", Status: "ativo", Capacity: 10}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/espacos/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Sala 1")
}

func TestGetAllSpacesHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "11111111111", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	_, _ = services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 1", Type: "sala", Status: "ativo", Capacity: 5})
	_, _ = services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 2", Type: "sala", Status: "ativo", Capacity: 8})

	router := setupSpaceRouter()

	req, _ := http.NewRequest(http.MethodGet, "/espacos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Sala 1")
	assert.Contains(t, resp.Body.String(), "Sala 2")
}

func TestGetSpaceByIDHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "22222222222", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	space, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 3", Type: "sala", Status: "ativo", Capacity: 10})

	router := setupSpaceRouter()

	req, _ := http.NewRequest(http.MethodGet, "/espacos/"+strconv.Itoa(int(space.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Sala 3")
}

func TestUpdateSpaceHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "33333333333", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	space, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 4", Type: "sala", Status: "ativo", Capacity: 10})

	router := setupSpaceRouter()

	update := dto.UpdateSpaceDTO{Name: "Sala Nova", Type: "laboratorio", Status: "inativo", Notice: "Fechada", Capacity: 20}
	jsonBody, _ := json.Marshal(update)

	req, _ := http.NewRequest(http.MethodPut, "/espacos/"+strconv.Itoa(int(space.ID)), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Sala Nova")
}

func TestDeleteSpaceHandler(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "44444444444", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	space, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 5", Type: "sala", Status: "ativo", Capacity: 10})

	router := setupSpaceRouter()

	req, _ := http.NewRequest(http.MethodDelete, "/espacos/"+strconv.Itoa(int(space.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestUpdateSpaceStatusHandler(t *testing.T) {
	tests.SetupTestDB(t)

	user, _ := services.CreateUser("User", "55555555555", "admin")
	token, _ := utils.GenerateJWT(user.ID, user.Role)

	space, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 6", Type: "sala", Status: "ativo", Capacity: 10})

	router := setupSpaceRouter()

	body := dto.UpdateStatusDTO{Status: "inativo"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPatch, "/espacos/"+strconv.Itoa(int(space.ID))+"/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Status atualizado com sucesso")
}

func TestUpdateSpaceNoticeHandler(t *testing.T) {
	tests.SetupTestDB(t)

	user, _ := services.CreateUser("User", "66666666666", "admin")
	token, _ := utils.GenerateJWT(user.ID, user.Role)

	space, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 7", Type: "sala", Status: "ativo", Capacity: 10})

	router := setupSpaceRouter()

	body := dto.UpdateNoticeDTO{Notice: "Aviso"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPatch, "/espacos/"+strconv.Itoa(int(space.ID))+"/aviso", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Aviso atualizado com sucesso")
}
