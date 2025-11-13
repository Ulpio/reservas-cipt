package config

import (
	"os"
	"strconv"
)

// EmailConfig contém as configurações SMTP
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Enabled  bool // Se o envio de email está habilitado
}

// GetEmailConfig retorna a configuração de email das variáveis de ambiente
func GetEmailConfig() EmailConfig {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	enabled, _ := strconv.ParseBool(getEnv("EMAIL_ENABLED", "false"))

	return EmailConfig{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     port,
		Username: getEnv("SMTP_USER", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("EMAIL_FROM", "noreply@cipt.com"),
		Enabled:  enabled,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

