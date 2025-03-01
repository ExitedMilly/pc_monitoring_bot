package main

import (
	"TG_BOT_GO/internal/config"
	"TG_BOT_GO/internal/telegram"
	"log"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	// Запускаем бота
	if err := telegram.StartBot(); err != nil {
		log.Fatalf("Ошибка при запуске бота: %v", err)
	}
}
