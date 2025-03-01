package monitor

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

// GetCPUUsage возвращает информацию о загруженности процессора в виде строки
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

// GetCPUUsageValue возвращает загрузку CPU в процентах (float64)
func GetCPUUsageValue() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		return 0.0
	}
	return percent[0]
}

// GetCPUTempValue возвращает температуру CPU в °C (float64)
func GetCPUTempValue() float64 {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0.0
	}

	for _, temp := range temps {
		if temp.SensorKey == "coretemp" || strings.Contains(temp.SensorKey, "CPU") {
			return temp.Temperature
		}
	}
	return 0.0
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
