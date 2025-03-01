package telegram

import (
	"TG_BOT_GO/internal/functions"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleUpdate обрабатывает входящие сообщения
func HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}

	switch update.Message.Command() {
	case "status":
		functions.HandleStatusCommand(update, bot)
	case "net":
		functions.HandleNetCommand(update, bot)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
		bot.Send(msg)
	}
}
