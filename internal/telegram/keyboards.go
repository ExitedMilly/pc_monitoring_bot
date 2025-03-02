package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// GetWelcomeKeyboard возвращает клавиатуру для приветственного сообщения
func GetWelcomeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мониторинг", "monitoring"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Управление", "management"),
		),
	)
}

// GetMonitoringKeyboard возвращает клавиатуру для раздела мониторинга
func GetMonitoringKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌐 Сетевой анализ", "network_analysis"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🖥️ Процессы", "processes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📈 Актуальное состояние", "status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚨 Предупреждения", "alarm"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back"),
		),
	)
}

// GetBackKeyboard возвращает клавиатуру с кнопкой "Назад"
func GetBackKeyboard(backTo string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", backTo),
		),
	)
}

// GetAlarmKeyboard возвращает клавиатуру для раздела предупреждений
func GetAlarmKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔔 Включить уведомления", "enable_alarm"),
			tgbotapi.NewInlineKeyboardButtonData("🔕 Выключить уведомления", "disable_alarm"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "monitoring"),
		),
	)
}
