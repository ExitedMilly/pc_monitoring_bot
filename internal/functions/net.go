package functions

import (
	"TG_BOT_GO/internal/monitor"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleNetCommand обрабатывает команду /net
func HandleNetCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Получаем информацию о IP
	ipInfo, err := monitor.GetIPInfo()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении информации о IP")
		bot.Send(msg)
		return
	}

	// Получаем топ процессов
	topProcesses, err := monitor.GetTopProcesses()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении информации о процессах")
		bot.Send(msg)
		return
	}

	// Получаем текущую скорость сети
	downloadSpeed, uploadSpeed, err := monitor.GetNetworkSpeed()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении скорости сети")
		bot.Send(msg)
		return
	}

	// Получаем общий трафик за 5 минут
	trafficLast5Min, err := monitor.GetTrafficLast5Min()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении трафика за 5 минут")
		bot.Send(msg)
		return
	}

	// Формируем вывод
	output := "+------------------------------+\n"
	output += "| 🌐 Сеть:                      \n"
	output += "+------------------------------+\n"
	output += fmt.Sprintf("🌍 Ваш IP: %s\n", ipInfo.IP)
	output += fmt.Sprintf("📍 Локация: %s, %s\n", ipInfo.City, ipInfo.Country)
	output += fmt.Sprintf("🏢 Провайдер: %s\n", ipInfo.Org)
	output += "\n"
	output += "📶 Текущая скорость:\n"
	output += fmt.Sprintf("  ⬇️ Входящая: %.2f МБ/с\n", downloadSpeed)
	output += fmt.Sprintf("  ⬆️ Исходящая: %.2f МБ/с\n", uploadSpeed)
	output += "\n"
	output += "📊 Топ процессов:\n"
	for i, p := range topProcesses {
		output += fmt.Sprintf("  %d. %s: ⬇️ %.1f МБ, ⬆️ %.1f МБ\n", i+1, p.Name, p.DownloadMB, p.UploadMB)
	}
	output += "\n"
	output += "📁 Общий трафик за 5 минут:\n"
	output += fmt.Sprintf("  ⬇️ Входящий: %.1f МБ\n", trafficLast5Min.DownloadMB)
	output += fmt.Sprintf("  ⬆️ Исходящий: %.1f МБ\n", trafficLast5Min.UploadMB)
	output += "\n"
	output += "📌 Топ-3 приложения за всё время:\n"
	for i, p := range topProcesses {
		output += fmt.Sprintf("  %d. %s: ⬇️ %.1f МБ, ⬆️ %.1f МБ\n", i+1, p.Name, p.DownloadMB, p.UploadMB)
	}
	output += "+------------------------------+"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
	bot.Send(msg)
}
