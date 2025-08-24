package services_test

import (
	"testing"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndRevokeStrike(t *testing.T) {
	tests.SetupTestDB(t)

	client := models.Client{Name: "Cliente", CPF: "123"}
	assert.NoError(t, database.DB.Create(&client).Error)

	input := dto.StrikeInputDTO{ClientID: client.ID, Reason: "Atraso"}
	created, err := services.CreateStrike(input)
	assert.NoError(t, err)
	assert.Equal(t, client.ID, created.ClientID)
	assert.False(t, created.Revoked)

	strikes, err := services.GetStrikesByClient(client.ID)
	assert.NoError(t, err)
	assert.Len(t, strikes, 1)

	// check client strike count
	var updated models.Client
	assert.NoError(t, database.DB.First(&updated, client.ID).Error)
	assert.Equal(t, 1, updated.Strikes)

	revoked, err := services.RevokeStrike(created.ID)
	assert.NoError(t, err)
	assert.True(t, revoked.Revoked)

	assert.NoError(t, database.DB.First(&updated, client.ID).Error)
	assert.Equal(t, 0, updated.Strikes)
}
