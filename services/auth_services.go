package services

import (
	"errors"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/Ulpio/reservas-cipt/utils"
)

func AuthenticateUser(cpf, password string) (string, error) {
	var user models.User
	if err := database.DB.Where("cpf = ?", cpf).First(&user).Error; err != nil {
		return "", errors.New("usuario nao encontrado")
	}
	if ok := utils.CheckPasswordHash(user.Password, password); !ok {
		return "", errors.New("senha incorreta")
	}
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", errors.New("erro ao gerar jwt")
	}
	return token, nil
}
