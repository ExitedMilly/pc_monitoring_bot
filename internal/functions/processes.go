package functions

import (
	"TG_BOT_GO/internal/monitor"
	"fmt"
	"sort"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shirou/gopsutil/process"
)

// ProcessInfo содержит информацию о процессе
type ProcessInfo struct {
	Name       string  // Название процесса
	CPUUsage   float64 // Нагрузка на CPU (%)
	GPULoad    float64 // Нагрузка на GPU (%)
	MemoryMB   float64 // Использование памяти (МБ)
	DownloadMB float64 // Входящий трафик (МБ)
	UploadMB   float64 // Исходящий трафик (МБ)
}

// HandleProcessesCommand обрабатывает команду /processes
func HandleProcessesCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Отправляем сообщение "Пожалуйста, подождите..."
	waitMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, подождите пару секунд...")
	sentMsg, _ := bot.Send(waitMsg)

	// Получаем список процессов
	processes, err := process.Processes()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении списка процессов")
		bot.Send(msg)
		return
	}

	// Собираем информацию о каждом процессе
	var processInfoList []ProcessInfo
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue // Пропускаем процесс, если не удалось получить имя
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue // Пропускаем процесс, если не удалось получить нагрузку на CPU
		}

		memInfo, err := p.MemoryInfo()
		if err != nil || memInfo == nil {
			continue // Пропускаем процесс, если не удалось получить информацию о памяти
		}

		gpuLoad := monitor.GetGPUUsageForProcess(p.Pid) // Нагрузка на GPU
		networkInfo, err := monitor.GetNetworkUsageForProcess(p.Pid)
		if err != nil {
			continue // Пропускаем процесс, если не удалось получить сетевую активность
		}

		processInfo := ProcessInfo{
			Name:       name,
			CPUUsage:   cpuPercent,
			GPULoad:    gpuLoad,
			MemoryMB:   float64(memInfo.RSS) / 1024 / 1024, // RSS в МБ
			DownloadMB: networkInfo.DownloadMB,
			UploadMB:   networkInfo.UploadMB,
		}
		processInfoList = append(processInfoList, processInfo)
	}

	// Сортируем процессы по убыванию нагрузки на CPU
	sort.Slice(processInfoList, func(i, j int) bool {
		return processInfoList[i].CPUUsage > processInfoList[j].CPUUsage
	})

	// Формируем вывод
	output := "+------------------------------+\n"
	output += "| 🖥️ Топ-10 процессов:          \n"
	output += "+------------------------------+\n"
	for i, p := range processInfoList {
		if i >= 10 { // Ограничиваемся 10 процессами
			break
		}
		output += fmt.Sprintf("%d. %s:\n", i+1, p.Name)
		output += fmt.Sprintf("  ⚙️ CPU: %.1f%%\n", p.CPUUsage)
		output += fmt.Sprintf("  🎮 GPU: %.1f%%\n", p.GPULoad)
		output += fmt.Sprintf("  🧠 Память: %.1f МБ\n", p.MemoryMB)
		output += fmt.Sprintf("  🌐 Сеть: ⬇️ %.1f МБ, ⬆️ %.1f МБ\n", p.DownloadMB, p.UploadMB)
		output += "\n" // Добавляем отступ между процессами
	}
	output += "+------------------------------+"

	// Удаляем сообщение "Пожалуйста, подождите..."
	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentMsg.MessageID)
	bot.Send(deleteMsg)

	// Отправляем результат
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
	bot.Send(msg)
}
