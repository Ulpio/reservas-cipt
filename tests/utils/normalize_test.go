package utils_test

import (
	"testing"

	"github.com/Ulpio/reservas-cipt/utils"
)

func TestNormalizeTime(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"16:00", "16:00:00", false},
		{"16:00:00", "16:00:00", false},
		{"09:30", "09:30:00", false},
		{"23:59:59", "23:59:59", false},
		{"16", "", true},           // Inválido
		{"4:00 PM", "", true},      // Inválido
		{"16:00:00:00", "", true},  // Inválido
	}

	for _, tt := range tests {
		result, err := utils.NormalizeTime(tt.input)
		
		if tt.hasError {
			if err == nil {
				t.Errorf("NormalizeTime(%q) deveria retornar erro, mas não retornou", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("NormalizeTime(%q) retornou erro inesperado: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("NormalizeTime(%q) = %q, esperado %q", tt.input, result, tt.expected)
			}
		}
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"2025-11-13", false},
		{"2025-01-01", false},
		{"13/11/2025", true},       // Formato brasileiro
		{"11-13-2025", true},       // Formato americano
		{"2025/11/13", true},       // Barras
		{"2025-13-01", true},       // Mês inválido
	}

	for _, tt := range tests {
		_, err := utils.ParseDate(tt.input)
		
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseDate(%q) deveria retornar erro, mas não retornou", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseDate(%q) retornou erro inesperado: %v", tt.input, err)
			}
		}
	}
}

func TestParseDateTime(t *testing.T) {
	result, err := utils.ParseDateTime("2025-11-13", "16:00")
	if err != nil {
		t.Errorf("ParseDateTime falhou: %v", err)
	}
	
	expectedHour := 16
	if result.Hour() != expectedHour {
		t.Errorf("Hora esperada %d, recebido %d", expectedHour, result.Hour())
	}

	expectedDay := 13
	if result.Day() != expectedDay {
		t.Errorf("Dia esperado %d, recebido %d", expectedDay, result.Day())
	}
}

