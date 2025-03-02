package functions

import (
	"TG_BOT_GO/internal/monitor"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleAlarmCommandOutput возвращает результат команды /alarm в виде строки
func HandleAlarmCommandOutput() string {
	output := "🚨 Уведомления: "
	if monitor.AlarmEnabled {
		output += "Включены\n"
	} else {
		output += "Выключены\n"
	}

	// Проверяем, установлены ли пороговые значения
	thresholdsSet := false
	if monitor.AlarmThresholds.CPUTemp > 0 {
		output += fmt.Sprintf("🌡️ Температура CPU: порог %.1f°C\n", monitor.AlarmThresholds.CPUTemp)
		thresholdsSet = true
	}
	if monitor.AlarmThresholds.GPUTemp > 0 {
		output += fmt.Sprintf("🌡️ Температура GPU: порог %.1f°C\n", monitor.AlarmThresholds.GPUTemp)
		thresholdsSet = true
	}
	if monitor.AlarmThresholds.CPUUsage > 0 {
		output += fmt.Sprintf("⚙️ Нагрузка CPU: порог %.1f%%\n", monitor.AlarmThresholds.CPUUsage)
		thresholdsSet = true
	}
	if monitor.AlarmThresholds.GPUUsage > 0 {
		output += fmt.Sprintf("🎮 Нагрузка GPU: порог %.1f%%\n", monitor.AlarmThresholds.GPUUsage)
		thresholdsSet = true
	}
	if monitor.AlarmThresholds.MemoryUsage > 0 {
		output += fmt.Sprintf("🧠 Память: порог %.1f%%\n", monitor.AlarmThresholds.MemoryUsage)
		thresholdsSet = true
	}
	if monitor.AlarmThresholds.NetworkUsage > 0 {
		output += fmt.Sprintf("🌐 Сеть: порог %.1f МБ/с\n", monitor.AlarmThresholds.NetworkUsage)
		thresholdsSet = true
	}
	if monitor.AlarmThresholds.DiskUsage > 0 {
		output += fmt.Sprintf("💽 Диски: порог %.1f%%\n", monitor.AlarmThresholds.DiskUsage)
		thresholdsSet = true
	}

	if !thresholdsSet {
		output += "\n⚠️ Ни одно пороговое значение не установлено. Используйте /alarm_set для настройки."
	}

	return output
}

// HandleAlarmCommand обрабатывает команду /alarm
func HandleAlarmCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	output := HandleAlarmCommandOutput()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
	bot.Send(msg)
}

// HandleAlarmOnCommand обрабатывает команду /alarm_on
func HandleAlarmOnCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if monitor.AlarmThresholds.CPUTemp == 0 && monitor.AlarmThresholds.GPUTemp == 0 &&
		monitor.AlarmThresholds.CPUUsage == 0 && monitor.AlarmThresholds.GPUUsage == 0 &&
		monitor.AlarmThresholds.MemoryUsage == 0 && monitor.AlarmThresholds.NetworkUsage == 0 &&
		monitor.AlarmThresholds.DiskUsage == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нельзя включить уведомления: пороговые значения не заданы.")
		bot.Send(msg)
		return
	}

	monitor.AlarmEnabled = true
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚨 Уведомления включены.")
	bot.Send(msg)
}

// HandleAlarmOffCommand обрабатывает команду /alarm_off
func HandleAlarmOffCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	monitor.AlarmEnabled = false
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚨 Уведомления выключены.")
	bot.Send(msg)
}

// HandleAlarmSetCommand обрабатывает команду /alarm_set
func HandleAlarmSetCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Использование: /alarm_set <параметр> <значение>")
		bot.Send(msg)
		return
	}

	param := args[0]
	value, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректное значение.")
		bot.Send(msg)
		return
	}

	switch param {
	case "cpu_temp":
		monitor.AlarmThresholds.CPUTemp = value
	case "gpu_tmp":
		monitor.AlarmThresholds.GPUTemp = value
	case "cpu_usage":
		monitor.AlarmThresholds.CPUUsage = value
	case "gpu_usage":
		monitor.AlarmThresholds.GPUUsage = value
	case "memory_usage":
		monitor.AlarmThresholds.MemoryUsage = value
	case "network_usage":
		monitor.AlarmThresholds.NetworkUsage = value
	case "disk_usage":
		monitor.AlarmThresholds.DiskUsage = value
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный параметр.")
		bot.Send(msg)
		return
	}

	if err := monitor.SaveThresholds(); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при сохранении пороговых значений.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("🚨 Порог для %s установлен на %.1f.", param, value))
	bot.Send(msg)
}
