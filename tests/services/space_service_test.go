package services_test

import (
	"testing"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetSpace(t *testing.T) {
	tests.SetupTestDB(t)

	input := dto.CreateSpaceDTO{Name: "Sala A", Type: "sala", Status: "ativo", Notice: "", Capacity: 10}
	created, err := services.CreateSpace(input)
	assert.NoError(t, err)
	assert.Equal(t, "Sala A", created.Name)

	fetched, err := services.GetSpaceByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "Sala A", fetched.Name)
}

func TestGetAllSpaces(t *testing.T) {
	tests.SetupTestDB(t)

	_, _ = services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 1", Type: "sala", Status: "ativo", Capacity: 5})
	_, _ = services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala 2", Type: "sala", Status: "ativo", Capacity: 8})

	spaces, err := services.GetAllSpaces()
	assert.NoError(t, err)
	assert.Len(t, spaces, 2)
}

func TestUpdateSpace(t *testing.T) {
	tests.SetupTestDB(t)

	created, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala", Type: "sala", Status: "ativo", Capacity: 10})

	update := dto.UpdateSpaceDTO{Name: "Sala X", Type: "laboratorio", Status: "inativo", Notice: "Fechada", Capacity: 20}
	updated, err := services.UpdateSpace(created.ID, update)
	assert.NoError(t, err)
	assert.Equal(t, "Sala X", updated.Name)
	assert.Equal(t, "laboratorio", updated.Type)
	assert.Equal(t, "inativo", updated.Status)
	assert.Equal(t, "Fechada", updated.Notice)
	assert.Equal(t, 20, updated.Capacity)
}

func TestDeleteSpace(t *testing.T) {
	tests.SetupTestDB(t)

	created, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala", Type: "sala", Status: "ativo", Capacity: 10})

	err := services.DeleteSpace(created.ID)
	assert.NoError(t, err)

	spaces, _ := services.GetAllSpaces()
	assert.Len(t, spaces, 0)
}

func TestUpdateSpaceStatus(t *testing.T) {
	tests.SetupTestDB(t)

	created, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala", Type: "sala", Status: "ativo", Capacity: 10})
	err := services.UpdateSpaceStatus(created.ID, "inativo")
	assert.NoError(t, err)

	updated, _ := services.GetSpaceByID(created.ID)
	assert.Equal(t, "inativo", updated.Status)
}

func TestUpdateSpaceNotice(t *testing.T) {
	tests.SetupTestDB(t)

	created, _ := services.CreateSpace(dto.CreateSpaceDTO{Name: "Sala", Type: "sala", Status: "ativo", Capacity: 10})
	err := services.UpdateSpaceNotice(created.ID, "Aviso")
	assert.NoError(t, err)

	updated, _ := services.GetSpaceByID(created.ID)
	assert.Equal(t, "Aviso", updated.Notice)
}
