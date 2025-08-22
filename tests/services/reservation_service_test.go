package services_test

import (
	"testing"
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndRetrieveReservation(t *testing.T) {
	tests.SetupTestDB(t)

	client := models.Client{Name: "Cliente", CPF: "123", Email: "cli@example.com"}
	assert.NoError(t, database.DB.Create(&client).Error)
	receptionist := models.User{Name: "Rec", CPF: "111", Role: "recepcionista"}
	assert.NoError(t, database.DB.Create(&receptionist).Error)
	space := models.Space{Name: "Sala", Type: "sala", Status: "ativo", Capacity: 10}
	assert.NoError(t, database.DB.Create(&space).Error)

	date := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	start := time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC)
	input := dto.CreateReservationDTO{
		ClientID:       client.ID,
		ReceptionistID: receptionist.ID,
		SpaceID:        space.ID,
		Date:           date,
		StartTime:      start,
		DurationHours:  2,
	}

	created, err := services.CreateReservation(input)
	assert.NoError(t, err)
	assert.Equal(t, client.ID, created.ClientID)

	fetched, err := services.GetReservationByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, 2, fetched.DurationHours)

	all, err := services.GetAllReservations()
	assert.NoError(t, err)
	assert.Len(t, all, 1)
}
