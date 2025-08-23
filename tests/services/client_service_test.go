package services_test

import (
	"testing"

	"github.com/Ulpio/reservas-cipt/dto"
	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/stretchr/testify/assert"
)

func TestBuscarOuCriarCliente(t *testing.T) {
	tests.SetupTestDB(t)

	input := dto.ClienteInputDTO{Name: "Joao", CPF: "12345678900", Email: "j@example.com", Phone: "999"}
	created, err := services.BuscarOuCriarCliente(input)
	assert.NoError(t, err)

	fetched, err := services.BuscarOuCriarCliente(input)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestGetClientByCPF(t *testing.T) {
	tests.SetupTestDB(t)

	_, err := services.GetClientByCPF("123")
	assert.Error(t, err)

	input := dto.ClienteInputDTO{Name: "Joao", CPF: "123", Email: "j@example.com", Phone: "999"}
	created, _ := services.BuscarOuCriarCliente(input)

	found, err := services.GetClientByCPF("123")
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestGetAllClientes(t *testing.T) {
	tests.SetupTestDB(t)

	_, _ = services.BuscarOuCriarCliente(dto.ClienteInputDTO{Name: "A", CPF: "1"})
	_, _ = services.BuscarOuCriarCliente(dto.ClienteInputDTO{Name: "B", CPF: "2"})

	clientes, err := services.GetAllClientes()
	assert.NoError(t, err)
	assert.Len(t, clientes, 2)
}
