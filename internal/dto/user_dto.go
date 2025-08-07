package dto

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	CPF      string `json:"cpf" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserIO struct {
	Name string `json:"name" binding:"required"`
	CPF  string `json:"cpf" binding:"required"`
	Role string `json:"role" binding:"required"`
}
