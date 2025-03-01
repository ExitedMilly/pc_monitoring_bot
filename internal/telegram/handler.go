package telegram

import (
	"TG_BOT_GO/internal/config"
	"TG_BOT_GO/internal/monitor"

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
		if update.Message != nil { // Игнорируем всё, кроме сообщений
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Обработка команды /start
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для мониторинга системы. Используй /status для проверки состояния.")
					bot.Send(msg)
				case "status":
					status := monitor.GetSystemStatus()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, status)
					bot.Send(msg)
				}
				continue
			}

			// Ответ на обычное сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я твой бот для мониторинга системы.")
			bot.Send(msg)
		}
	}

	return nil
}
