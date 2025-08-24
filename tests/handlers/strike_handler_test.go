package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	"github.com/stretchr/testify/assert"
)

func setupStrikeRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("")
	strikes := api.Group("/strikes")
	strikes.Use(middleware.JWTAuthMiddleware())
	strikes.POST("", handlers.CreateStrikeHandler)
	strikes.GET("/client/:id", handlers.GetStrikesByClientHandler)
	strikes.DELETE("/:id", handlers.RevokeStrikeHandler)
	return r
}

func TestStrikeHandlers(t *testing.T) {
	tests.SetupTestDB(t)

	admin, _ := services.CreateUser("Admin", "00000000000", "admin")
	token, _ := utils.GenerateJWT(admin.ID, admin.Role)

	client := models.Client{Name: "Cliente", CPF: "123"}
	assert.NoError(t, database.DB.Create(&client).Error)

	router := setupStrikeRouter()

	body := dto.StrikeInputDTO{ClientID: client.ID, Reason: "Atraso"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/strikes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var created dto.StrikeOutputDTO
	json.Unmarshal(resp.Body.Bytes(), &created)

	// list
	reqGet, _ := http.NewRequest(http.MethodGet, "/strikes/client/"+strconv.Itoa(int(client.ID)), nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	respGet := httptest.NewRecorder()
	router.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)
	assert.Contains(t, respGet.Body.String(), "Atraso")

	// revoke
	reqDel, _ := http.NewRequest(http.MethodDelete, "/strikes/"+strconv.Itoa(int(created.ID)), nil)
	reqDel.Header.Set("Authorization", "Bearer "+token)
	respDel := httptest.NewRecorder()
	router.ServeHTTP(respDel, reqDel)
	assert.Equal(t, http.StatusOK, respDel.Code)
	assert.Contains(t, respDel.Body.String(), "true")
}
