package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não foi carregado automaticamente:", err)
	} else {
		log.Println("✅ .env carregado com sucesso")
	}
}
