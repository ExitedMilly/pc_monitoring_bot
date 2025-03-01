package monitor

import (
	"fmt"
	"os/exec"
	"runtime"
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
	gpuUsage := getGPUUsage() // Переместили видеокарту после процессора
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

// Получение информации о загруженности видеокарты
func getGPUUsage() string {
	gpuType := detectGPU()

	switch gpuType {
	case "nvidia":
		usage, temp, err := getNvidiaGPUUsage()
		if err == nil {
			return fmt.Sprintf("NVIDIA:\n🔄 Загруженность GPU: %s%%\n%s\n🌡️ Температура: %.1f°C", usage, getProgressBar(parseFloat(usage)), temp)
		}
	case "amd":
		usage, temp, err := getAMDGPUUsage()
		if err == nil {
			return fmt.Sprintf("AMD:\n🔄 Загруженность GPU: %s%%\n%s\n🌡️ Температура: %.1f°C", usage, getProgressBar(parseFloat(usage)), temp)
		}
	case "intel":
		usage, temp, err := getIntelGPUUsage()
		if err == nil {
			return fmt.Sprintf("Intel:\n🔄 Загруженность GPU: %s%%\n%s\n🌡️ Температура: %.1f°C", usage, getProgressBar(parseFloat(usage)), temp)
		}
	}

	return "Видеокарта: не удалось получить данные"
}

// Получение загруженности и температуры видеокарты NVIDIA
func getNvidiaGPUUsage() (string, float64, error) {
	out, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu,temperature.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return "", 0, err
	}
	data := strings.Split(strings.TrimSpace(string(out)), ", ")
	if len(data) < 2 {
		return "", 0, fmt.Errorf("неверный формат данных")
	}
	usage := data[0]
	temp := parseFloat(data[1])
	return usage, temp, nil
}

// Получение загруженности и температуры видеокарты AMD
func getAMDGPUUsage() (string, float64, error) {
	out, err := exec.Command("rocm-smi", "--showuse", "--showtemp").Output()
	if err != nil {
		return "", 0, err
	}
	lines := strings.Split(string(out), "\n")
	var usage, temp string
	for _, line := range lines {
		if strings.Contains(line, "GPU use") {
			usage = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Temperature") {
			temp = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return usage, parseFloat(temp), nil
}

// Получение загруженности и температуры видеокарты Intel
func getIntelGPUUsage() (string, float64, error) {
	out, err := exec.Command("intel_gpu_top", "-l", "-o", "-").Output()
	if err != nil {
		return "", 0, err
	}
	lines := strings.Split(string(out), "\n")
	var usage, temp string
	for _, line := range lines {
		if strings.Contains(line, "Render/3D") {
			usage = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Temperature") {
			temp = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return usage, parseFloat(temp), nil
}

// Преобразование строки в число
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
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

// Определение типа видеокарты
func detectGPU() string {
	switch runtime.GOOS {
	case "linux":
		return detectGPULinux()
	case "windows":
		return detectGPUWindows()
	case "darwin":
		return detectGPUMacOS()
	default:
		return "unknown"
	}
}

// Определение видеокарты на Linux
func detectGPULinux() string {
	out, err := exec.Command("lspci").Output()
	if err != nil {
		return "unknown"
	}
	if strings.Contains(strings.ToLower(string(out)), "nvidia") {
		return "nvidia"
	}
	if strings.Contains(strings.ToLower(string(out)), "amd") {
		return "amd"
	}
	if strings.Contains(strings.ToLower(string(out)), "intel") {
		return "intel"
	}
	return "unknown"
}

// Определение видеокарты на Windows
func detectGPUWindows() string {
	out, err := exec.Command("wmic", "path", "win32_VideoController", "get", "name").Output()
	if err != nil {
		return "unknown"
	}
	if strings.Contains(strings.ToLower(string(out)), "nvidia") {
		return "nvidia"
	}
	if strings.Contains(strings.ToLower(string(out)), "amd") {
		return "amd"
	}
	if strings.Contains(strings.ToLower(string(out)), "intel") {
		return "intel"
	}
	return "unknown"
}

// Определение видеокарты на macOS
func detectGPUMacOS() string {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err != nil {
		return "unknown"
	}
	if strings.Contains(strings.ToLower(string(out)), "nvidia") {
		return "nvidia"
	}
	if strings.Contains(strings.ToLower(string(out)), "amd") {
		return "amd"
	}
	if strings.Contains(strings.ToLower(string(out)), "intel") {
		return "intel"
	}
	return "unknown"
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
		if temp.SensorKey == "coretemp" || strings.Contains(temp.SensorKey, "CPU") {
			tempInfo += fmt.Sprintf("🌡️ Температура: %.1f°C\n", temp.Temperature)
		}
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
