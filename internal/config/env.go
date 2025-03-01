package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig загружает переменные окружения из .env файла
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
}

// GetEnv возвращает значение переменной окружения
func GetEnv(key string) string {
	return os.Getenv(key)
}
