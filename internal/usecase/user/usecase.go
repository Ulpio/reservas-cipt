package user

import (
	"errors"
	"strings"

	"github.com/Ulpio/reservas-cipt/internal/domain/user"
)

type UserUseCase struct {
	repo user.Repository
}

func NewUserUseCase(repo user.Repository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) RegisterUser(nomeCompleto, cpf, role, password string) error {
	if strings.TrimSpace(nomeCompleto) == "" {
		return errors.New("nome completo é obrigatório")
	}

	existingUser, err := uc.repo.FindByCPF(cpf)
	if err == nil && existingUser != nil {
		return errors.New("CPF já cadastrado")
	}

	hashedPassword := password

	newUser := &user.User{
		Name:     nomeCompleto,
		CPF:      cpf,
		Role:     role,
		Password: hashedPassword,
	}

	return uc.repo.Create(newUser)
}

func (uc *UserUseCase) FindUserByID(id uint) (*user.User, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) FindAllUsers() ([]user.User, error) {
	return uc.repo.FindAll()
}

func (uc *UserUseCase) DeleteUser(id uint) error {
	err := uc.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) FindUserByCPF(cpf string) (*user.User, error) {
	user, err := uc.repo.FindByCPF(cpf)
	if err != nil {
		return nil, err
	}
	return user, nil
}
