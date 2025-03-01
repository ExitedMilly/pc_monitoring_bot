package monitor

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
)

// GetDiskUsage возвращает информацию о дисках в виде строки
func GetDiskUsage() string {
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

// GetDiskUsageValue возвращает загруженность дисков в процентах (float64)
func GetDiskUsageValue() float64 {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0.0
	}
	return usage.UsedPercent
}
