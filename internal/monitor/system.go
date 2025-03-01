package monitor

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// GetSystemStatus возвращает строку с информацией о системе
func GetSystemStatus() string {
	diskInfo := getDiskUsage()
	cpuUsage := getCPUUsage()
	memInfo := getMemoryUsage()
	networkInfo := getNetworkUsage()

	// Форматируем вывод с рамками
	output := "+------------------------------+\n"
	output += "| 💽 Диски:                     \n"
	output += "+------------------------------+\n"
	output += diskInfo
	output += "+------------------------------+\n"
	output += "| ⚙️ Процессор:                 \n"
	output += "+------------------------------+\n"
	output += cpuUsage + "\n"
	output += "+------------------------------+\n"
	output += "| 🧠 Память:                    \n"
	output += "+------------------------------+\n"
	output += memInfo + "\n"
	output += "+------------------------------+\n"
	output += "| 🌐 Сеть:                      \n"
	output += "+------------------------------+\n"
	output += networkInfo + "\n"
	output += "+------------------------------+"

	return output
}

// Получение информации о всех дисках
func getDiskUsage() string {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return "Ошибка при получении информации о дисках"
	}

	var diskInfo string
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		freeGB := float64(usage.Free) / 1024 / 1024 / 1024
		totalGB := float64(usage.Total) / 1024 / 1024 / 1024
		diskInfo += fmt.Sprintf(
			"📁 Диск %s: %.1f ГБ свободно из %.1f ГБ (%.1f%%)\n",
			partition.Mountpoint,
			freeGB,
			totalGB,
			100-usage.UsedPercent,
		)
	}

	return diskInfo
}

// Получение информации о процессоре
func getCPUUsage() string {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return "Ошибка при получении информации о процессоре"
	}
	progressBar := getProgressBar(percent[0])

	// Получение температуры процессора
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return fmt.Sprintf("🔄 Загрузка: %.2f%%\n%s", percent[0], progressBar)
	}

	var tempInfo string
	for _, temp := range temps {
		if strings.Contains(temp.SensorKey, "coretemp") || strings.Contains(temp.SensorKey, "CPU") {
			tempInfo += fmt.Sprintf("🌡️ Температура: %.1f°C\n", temp.Temperature)
		}
	}

	if tempInfo == "" {
		tempInfo = "🌡️ Температура: недоступно\n"
	}

	return fmt.Sprintf("🔄 Загрузка: %.2f%%\n%s%s", percent[0], progressBar, tempInfo)
}

// Получение информации о памяти
func getMemoryUsage() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "Ошибка при получении информации о памяти"
	}

	usedPercent := memInfo.UsedPercent
	progressBar := getProgressBar(usedPercent)
	return fmt.Sprintf(
		"📊 Занято: %.1f ГБ из %.1f ГБ (%.2f%%)\n%s",
		float64(memInfo.Used)/1024/1024/1024,
		float64(memInfo.Total)/1024/1024/1024,
		usedPercent,
		progressBar,
	)
}

// Получение информации о сети
func getNetworkUsage() string {
	io, err := net.IOCounters(true)
	if err != nil {
		return "Ошибка при получении информации о сети"
	}

	var networkInfo string
	for _, stats := range io {
		// Пропускаем интерфейсы с нулевым трафиком
		if stats.BytesRecv == 0 && stats.BytesSent == 0 {
			continue
		}

		// Форматируем вывод
		networkInfo += fmt.Sprintf(
			"📶 Интерфейс %s:\n  ⬇️ Входящий: %.1f МБ\n  ⬆️ Исходящий: %.1f МБ\n",
			stats.Name,
			float64(stats.BytesRecv)/1024/1024,
			float64(stats.BytesSent)/1024/1024,
		)
	}

	return networkInfo
}

// getProgressBar возвращает строку с полоской загрузки
func getProgressBar(percent float64) string {
	const barLength = 10 // Длина полоски (10 символов)

	// Определяем количество заполненных символов
	filled := int((percent + 5) / 10) // Округляем до ближайшего десятка
	if filled > barLength {
		filled = barLength
	}

	// Создаём полоску
	bar := strings.Repeat("■", filled) + strings.Repeat("▢", barLength-filled)
	return bar
}
