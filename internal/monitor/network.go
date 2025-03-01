package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/net"
)

// GetNetworkUsage возвращает информацию о сети
func GetNetworkUsage() string {
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
