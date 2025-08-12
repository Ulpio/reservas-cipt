package dto

type LoginInputDTO struct {
	CPF      string `json:"cpf" binding:"required"`
	Password string `json:"password" binding:"required"`
}
