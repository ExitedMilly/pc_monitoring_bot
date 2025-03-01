package telegram

import (
	"TG_BOT_GO/internal/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartBot запускает телеграм бота
func StartBot() error {
	// Получаем токен бота из конфигурации
	botToken := config.GetEnv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен")
	}

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	// Включаем режим отладки (опционально)
	bot.Debug = true

	log.Printf("Бот %s успешно запущен", bot.Self.UserName)

	// Настраиваем канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	for update := range updates {
		HandleUpdate(update, bot)
	}

	return nil
}
