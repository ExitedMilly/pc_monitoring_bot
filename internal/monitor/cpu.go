package monitor

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

// GetCPUUsage возвращает информацию о загруженности процессора
func GetCPUUsage() string {
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
		if temp.SensorKey == "coretemp" || strings.Contains(temp.SensorKey, "CPU") {
			tempInfo += fmt.Sprintf("🌡️ Температура: %.1f°C\n", temp.Temperature)
		}
	}

	return fmt.Sprintf("🔄 Загрузка: %.2f%%\n%s%s", percent[0], progressBar, tempInfo)
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
