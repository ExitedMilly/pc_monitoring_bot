package telegram

import (
	"TG_BOT_GO/internal/functions"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleUpdate обрабатывает входящие сообщения и callback-запросы
func HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		// Обработка текстовых команд
		switch update.Message.Command() {
		case "start":
			sendWelcomeMessage(update.Message.Chat.ID, 0, bot) // messageID = 0 для нового сообщения
		case "net":
			functions.HandleNetCommand(update, bot)
		case "processes":
			functions.HandleProcessesCommand(update, bot)
		case "status":
			functions.HandleStatusCommand(update, bot)
		case "alarm":
			functions.HandleAlarmCommand(update, bot)
		case "alarm_on":
			functions.HandleAlarmOnCommand(update, bot)
		case "alarm_off":
			functions.HandleAlarmOffCommand(update, bot)
		case "alarm_set":
			functions.HandleAlarmSetCommand(update, bot)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
			bot.Send(msg)
		}
	} else if update.CallbackQuery != nil {
		// Обработка callback-запросов от inline-кнопок
		handleCallbackQuery(update.CallbackQuery, bot)
	}
}

// handleCallbackQuery обрабатывает нажатия на inline-кнопки
func handleCallbackQuery(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	data := callback.Data

	switch data {
	case "monitoring":
		sendMonitoringMenu(chatID, messageID, bot)
	case "network_analysis":
		handleNetworkAnalysis(chatID, messageID, bot)
	case "processes":
		handleProcesses(chatID, messageID, bot)
	case "status":
		handleStatus(chatID, messageID, bot)
	case "alarm":
		handleAlarm(chatID, messageID, bot)
	case "enable_alarm":
		handleEnableAlarm(chatID, messageID, bot)
	case "disable_alarm":
		handleDisableAlarm(chatID, messageID, bot)
	case "back":
		sendWelcomeMessage(chatID, messageID, bot) // Возврат в главное меню
	default:
		// Обработка других callback-запросов
	}
}

// handleEnableAlarm обрабатывает запрос на включение уведомлений
func handleEnableAlarm(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	// Выполняем команду /alarm_on
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}
	functions.HandleAlarmOnCommand(update, bot)

	// Редактируем сообщение на результат команды /alarm
	handleAlarm(chatID, messageID, bot)
}

// handleDisableAlarm обрабатывает запрос на выключение уведомлений
func handleDisableAlarm(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	// Выполняем команду /alarm_off
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}
	functions.HandleAlarmOffCommand(update, bot)

	// Редактируем сообщение на результат команды /alarm
	handleAlarm(chatID, messageID, bot)
}

// sendWelcomeMessage отправляет приветственное сообщение с inline-кнопками
func sendWelcomeMessage(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	text := `Привет! 👋 Я бот для удалённого мониторинга и управления компьютером!
Благодаря мне, ты сможешь улучшить свой опыт пользования компьютером, облегчить себе жизнь и отследить вредоносное ПО!
Выбери раздел, который тебя интересует.`

	keyboard := GetWelcomeKeyboard()

	// Если messageID == 0, отправляем новое сообщение
	if messageID == 0 {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = &keyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
		}
		return
	}

	// Редактируем существующее сообщение
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ReplyMarkup = &keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при редактировании сообщения: %v", err)
	}
}

// sendMonitoringMenu отправляет меню мониторинга
func sendMonitoringMenu(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	text := `📊 Раздел мониторинга позволяет анализировать интернет-трафик, работу комплектующих вашего ПК, просматривать процессы на компьютере и настраивать экстренные уведомления о чрезмерной температуре или нагрузке на комплектующие.`

	keyboard := GetMonitoringKeyboard()

	// Если messageID == 0, отправляем новое сообщение
	if messageID == 0 {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = &keyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
		}
		return
	}

	// Редактируем существующее сообщение
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ReplyMarkup = &keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при редактировании сообщения: %v", err)
	}
}

// handleNetworkAnalysis обрабатывает запрос на сетевой анализ
func handleNetworkAnalysis(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	// Отправляем сообщение "Пожалуйста, подождите пару секунд"
	msg := tgbotapi.NewEditMessageText(chatID, messageID, "⏳ Пожалуйста, подождите пару секунд...")
	bot.Send(msg)

	// Получаем результат работы /net
	output := functions.HandleNetCommandOutput()

	// Формируем результат
	resultMsg := tgbotapi.NewEditMessageText(chatID, messageID, output)
	keyboard := GetBackKeyboard("monitoring") // Кнопка "Назад"
	resultMsg.ReplyMarkup = &keyboard

	// Редактируем сообщение с результатом
	_, err := bot.Send(resultMsg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// handleProcesses обрабатывает запрос на просмотр процессов
func handleProcesses(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	// Отправляем сообщение "Пожалуйста, подождите пару секунд"
	msg := tgbotapi.NewEditMessageText(chatID, messageID, "⏳ Пожалуйста, подождите пару секунд...")
	bot.Send(msg)

	// Получаем результат работы /processes
	output := functions.HandleProcessesCommandOutput()

	// Формируем результат
	resultMsg := tgbotapi.NewEditMessageText(chatID, messageID, output)
	keyboard := GetBackKeyboard("monitoring") // Кнопка "Назад"
	resultMsg.ReplyMarkup = &keyboard

	// Редактируем сообщение с результатом
	_, err := bot.Send(resultMsg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// handleStatus обрабатывает запрос на актуальное состояние
func handleStatus(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	// Отправляем сообщение "Пожалуйста, подождите пару секунд"
	msg := tgbotapi.NewEditMessageText(chatID, messageID, "⏳ Пожалуйста, подождите пару секунд...")
	bot.Send(msg)

	// Получаем результат работы /status
	output := functions.HandleStatusCommandOutput()

	// Формируем результат
	resultMsg := tgbotapi.NewEditMessageText(chatID, messageID, output)
	keyboard := GetBackKeyboard("monitoring") // Кнопка "Назад"
	resultMsg.ReplyMarkup = &keyboard

	// Редактируем сообщение с результатом
	_, err := bot.Send(resultMsg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// handleAlarm обрабатывает запрос на настройку предупреждений
func handleAlarm(chatID int64, messageID int, bot *tgbotapi.BotAPI) {
	// Получаем результат работы /alarm
	output := functions.HandleAlarmCommandOutput()

	// Редактируем текущее сообщение
	msg := tgbotapi.NewEditMessageText(chatID, messageID, output)
	keyboard := GetAlarmKeyboard() // Используем клавиатуру для раздела предупреждений
	msg.ReplyMarkup = &keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при редактировании сообщения: %v", err)
	}
}
