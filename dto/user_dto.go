package dto

type UserInputDTO struct {
	Name string `json:"name"`
	CPF  string `json:"cpf"`
	Role string `json:"role"`
}

type UserOutputDTO struct {
	ID   uint   `json:"user_id"`
	Name string `json:"name"`
	CPF  string `json:"cpf"`
	Role string `json:"role"`
}
