package utils_test

import (
	"testing"
	"time"

	"github.com/Ulpio/reservas-cipt/utils"
)

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		cpf       string
		shouldErr bool
		desc      string
	}{
		{"11144477735", false, "CPF válido"},
		{"111.444.777-35", false, "CPF válido com formatação"},
		{"00000000000", true, "CPF com todos os dígitos iguais"},
		{"11111111111", true, "CPF com todos os dígitos iguais"},
		{"12345678901", true, "CPF com dígitos verificadores inválidos"},
		{"123456789", true, "CPF com menos de 11 dígitos"},
		{"123456789012", true, "CPF com mais de 11 dígitos"},
	}

	for _, tt := range tests {
		err := utils.ValidateCPF(tt.cpf)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("%s: esperava erro, mas não houve erro", tt.desc)
			}
		} else {
			if err != nil {
				t.Errorf("%s: não esperava erro, mas houve: %v", tt.desc, err)
			}
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email     string
		shouldErr bool
		desc      string
	}{
		{"", false, "Email vazio (opcional)"},
		{"joao@example.com", false, "Email válido"},
		{"maria.silva@empresa.com.br", false, "Email válido com múltiplos domínios"},
		{"invalido@", true, "Email sem domínio"},
		{"@example.com", true, "Email sem usuário"},
		{"invalido", true, "Email sem @"},
		{"invalido@.", true, "Email com domínio inválido"},
	}

	for _, tt := range tests {
		err := utils.ValidateEmail(tt.email)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("%s: esperava erro, mas não houve erro", tt.desc)
			}
		} else {
			if err != nil {
				t.Errorf("%s: não esperava erro, mas houve: %v", tt.desc, err)
			}
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		phone     string
		shouldErr bool
		desc      string
	}{
		{"", false, "Telefone vazio (opcional)"},
		{"11999887766", false, "Telefone válido com 11 dígitos"},
		{"1133334444", false, "Telefone válido com 10 dígitos"},
		{"(11) 99988-7766", false, "Telefone formatado"},
		{"119998877", true, "Telefone com menos de 10 dígitos"},
		{"119998877665", true, "Telefone com mais de 11 dígitos"},
	}

	for _, tt := range tests {
		err := utils.ValidatePhoneNumber(tt.phone)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("%s: esperava erro, mas não houve erro", tt.desc)
			}
		} else {
			if err != nil {
				t.Errorf("%s: não esperava erro, mas houve: %v", tt.desc, err)
			}
		}
	}
}

func TestValidateMinimumAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		birthDate time.Time
		minAge    int
		shouldErr bool
		desc      string
	}{
		{now.AddDate(-20, 0, 0), 18, false, "20 anos de idade (válido)"},
		{now.AddDate(-18, 0, -1), 18, false, "18 anos e 1 dia (válido)"},
		{now.AddDate(-17, 0, 0), 18, true, "17 anos exatos (inválido)"},
		{now.AddDate(-10, 0, 0), 18, true, "10 anos (inválido)"},
		{now.AddDate(-25, 0, 0), 21, false, "25 anos com mínimo de 21 (válido)"},
	}

	for _, tt := range tests {
		err := utils.ValidateMinimumAge(tt.birthDate, tt.minAge)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("%s: esperava erro, mas não houve erro", tt.desc)
			}
		} else {
			if err != nil {
				t.Errorf("%s: não esperava erro, mas houve: %v", tt.desc, err)
			}
		}
	}
}

