package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

// GetMemoryUsage возвращает информацию о памяти в виде строки
func GetMemoryUsage() string {
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

// GetMemoryUsageValue возвращает использование памяти в процентах (float64)
func GetMemoryUsageValue() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0.0
	}
	return memInfo.UsedPercent
}
