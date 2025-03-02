package functions

import (
	"TG_BOT_GO/internal/monitor"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleStatusCommandOutput возвращает результат команды /status в виде строки
func HandleStatusCommandOutput() string {
	diskInfo := monitor.GetDiskUsage()
	cpuUsage := monitor.GetCPUUsage()
	gpuUsage := monitor.GetGPUUsage()
	memInfo := monitor.GetMemoryUsage()
	networkInfo := monitor.GetNetworkUsage()

	output := "+------------------------------+\n"
	output += "| 💽 Диски:                     \n"
	output += "+------------------------------+\n"
	output += diskInfo
	output += "+------------------------------+\n"
	output += "| ⚙️ Процессор:                 \n"
	output += "+------------------------------+\n"
	output += cpuUsage + "\n"
	output += "+------------------------------+\n"
	output += "| 🎮 Видеокарта:                \n"
	output += "+------------------------------+\n"
	output += gpuUsage + "\n"
	output += "+------------------------------+\n"
	output += "| 🧠 Память:                    \n"
	output += "+------------------------------+\n"
	output += memInfo + "\n"
	output += "+------------------------------+\n"
	output += "| 🌐 Сеть:                      \n"
	output += "+------------------------------+\n"
	output += networkInfo
	output += "+------------------------------+"

	return output
}

// HandleStatusCommand обрабатывает команду /status
func HandleStatusCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	output := HandleStatusCommandOutput()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
	bot.Send(msg)
}
