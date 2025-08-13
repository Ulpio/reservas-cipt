package services_test

import (
	"testing"

	"github.com/Ulpio/reservas-cipt/services"
	"github.com/Ulpio/reservas-cipt/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetUser(t *testing.T) {
	tests.SetupTestDB(t)

	// Cria o usuário
	created, err := services.CreateUser("Alice", "12345678900", "admin")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", created.Name)

	// Busca pelo ID
	user, err := services.GetUserByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.CPF, user.CPF)
	assert.Equal(t, created.Role, user.Role)
}

func TestGetAllUsers(t *testing.T) {
	tests.SetupTestDB(t)

	// Cria múltiplos usuários
	_, _ = services.CreateUser("User 1", "11111111111", "admin")
	_, _ = services.CreateUser("User 2", "22222222222", "receptionist")

	users, err := services.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestDeleteUser(t *testing.T) {
	tests.SetupTestDB(t)

	u, _ := services.CreateUser("Temp", "99999999999", "admin")

	err := services.DeleteUser(u.ID)
	assert.NoError(t, err)

	_, err = services.GetUserByID(u.ID)
	assert.Error(t, err) // esperado: not found
}
