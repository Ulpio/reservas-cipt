package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NormalizeCPF remove pontos, traços e espaços do CPF
func NormalizeCPF(cpf string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(cpf, "")
}

// ValidateCPF valida se o CPF é válido (formato e dígitos verificadores)
func ValidateCPF(cpf string) error {
	// Normalizar CPF
	cpf = NormalizeCPF(cpf)
	
	// Verificar se tem 11 dígitos
	if len(cpf) != 11 {
		return errors.New("CPF deve ter 11 dígitos")
	}
	
	// Verificar se todos os dígitos são iguais (CPF inválido)
	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return errors.New("CPF inválido")
	}
	
	// Validar primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}
	
	expectedFirst, _ := strconv.Atoi(string(cpf[9]))
	if firstDigit != expectedFirst {
		return errors.New("CPF inválido")
	}
	
	// Validar segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}
	
	expectedSecond, _ := strconv.Atoi(string(cpf[10]))
	if secondDigit != expectedSecond {
		return errors.New("CPF inválido")
	}
	
	return nil
}

// ValidateEmail valida formato de email
func ValidateEmail(email string) error {
	if email == "" {
		return nil // Email é opcional
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("formato de email inválido")
	}
	
	return nil
}

// ValidatePhoneNumber valida formato de telefone
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return nil // Telefone é opcional
	}
	
	normalized := NormalizePhone(phone)
	length := len(normalized)
	
	if length < 10 || length > 11 {
		return errors.New("telefone deve ter entre 10 e 11 dígitos")
	}
	
	return nil
}

// ValidateMinimumAge valida se a pessoa tem idade mínima
func ValidateMinimumAge(birthDate time.Time, minimumAge int) error {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	
	// Ajustar se ainda não fez aniversário este ano
	if now.Month() < birthDate.Month() || 
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	
	if age < minimumAge {
		return fmt.Errorf("cliente deve ter pelo menos %d anos", minimumAge)
	}
	
	return nil
}

// NormalizePhone remove pontos, traços, parênteses e espaços do telefone
func NormalizePhone(phone string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(phone, "")
}

// NormalizeTime normaliza horário no formato HH:MM ou HH:MM:SS para HH:MM:SS
func NormalizeTime(timeStr string) (string, error) {
	timeStr = strings.TrimSpace(timeStr)
	
	// Aceita HH:MM:SS
	if regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`).MatchString(timeStr) {
		return timeStr, nil
	}
	
	// Aceita HH:MM e adiciona :00
	if regexp.MustCompile(`^\d{2}:\d{2}$`).MatchString(timeStr) {
		return timeStr + ":00", nil
	}
	
	return "", errors.New("formato de hora inválido. Use HH:MM ou HH:MM:SS")
}

// ParseDate parseia data no formato YYYY-MM-DD para time.Time
func ParseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	
	// Parsear no formato YYYY-MM-DD usando localização local
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return time.Time{}, errors.New("formato de data inválido. Use YYYY-MM-DD")
	}
	
	return date, nil
}

// ParseDateTime parseia data e hora separadas e retorna time.Time completo
func ParseDateTime(dateStr, timeStr string) (time.Time, error) {
	// Normalizar hora
	normalizedTime, err := NormalizeTime(timeStr)
	if err != nil {
		return time.Time{}, err
	}
	
	// Parsear data
	date, err := ParseDate(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	
	// Combinar data e hora usando a localização local do servidor
	dateTimeStr := fmt.Sprintf("%s %s", date.Format("2006-01-02"), normalizedTime)
	// Usar Local ao invés de UTC
	dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateTimeStr, time.Local)
	if err != nil {
		return time.Time{}, errors.New("erro ao combinar data e hora")
	}
	
	return dateTime, nil
}

// ValidateFutureDate valida se a data não é passada
func ValidateFutureDate(date time.Time) error {
	now := time.Now()
	// Usar o mesmo location para ambas as datas
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, now.Location())
	
	if dateOnly.Before(today) {
		return errors.New("a data não pode ser no passado")
	}
	
	return nil
}

