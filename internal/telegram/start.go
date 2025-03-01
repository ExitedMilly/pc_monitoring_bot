package telegram

import (
	"TG_BOT_GO/internal/config"
	"TG_BOT_GO/internal/monitor"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	chatID int64 // Глобальная переменная для хранения chatID
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

	// Загружаем пороговые значения
	if err := monitor.LoadThresholds(); err != nil {
		log.Println("Ошибка при загрузке пороговых значений:", err)
	}

	// Запускаем мониторинг уведомлений
	go monitor.StartAlarmMonitor(bot, chatID)

	// Настраиваем канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	for update := range updates {
		if update.Message != nil {
			// Сохраняем chatID при первом взаимодействии
			if chatID == 0 {
				chatID = update.Message.Chat.ID
				log.Printf("ChatID сохранен: %d", chatID)

				// Перезапускаем мониторинг уведомлений с новым chatID
				go monitor.StartAlarmMonitor(bot, chatID)
			}
		}
		HandleUpdate(update, bot)
	}

	return nil
}
