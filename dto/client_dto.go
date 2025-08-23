package dto

type ClienteInputDTO struct {
	Name  string `json:"name"`
	CPF   string `json:"cpf" binding:"required"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ClienteOutputDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	CPF     string `json:"cpf"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Strikes int    `json:"strikes"`
}
