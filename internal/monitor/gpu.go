package monitor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// GetGPUUsage возвращает информацию о видеокарте
func GetGPUUsage() string {
	gpuType := detectGPU()

	switch gpuType {
	case "nvidia":
		usage, temp, err := getNvidiaGPUUsage()
		if err == nil {
			return fmt.Sprintf("🎮 NVIDIA:\n🔄 Загруженность GPU: %s%%\n%s\n🌡️ Температура: %.1f°C", usage, getProgressBar(parseFloat(usage)), temp)
		}
	case "amd":
		usage, temp, err := getAMDGPUUsage()
		if err == nil {
			return fmt.Sprintf("🎮 AMD:\n🔄 Загруженность GPU: %s%%\n%s\n🌡️ Температура: %.1f°C", usage, getProgressBar(parseFloat(usage)), temp)
		}
	case "intel":
		usage, temp, err := getIntelGPUUsage()
		if err == nil {
			return fmt.Sprintf("🎮 Intel:\n🔄 Загруженность GPU: %s%%\n%s\n🌡️ Температура: %.1f°C", usage, getProgressBar(parseFloat(usage)), temp)
		}
	}

	return "🎮 Видеокарта: не удалось получить данные"
}

// detectGPU определяет тип видеокарты
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

// detectGPULinux определяет видеокарту на Linux
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

// detectGPUWindows определяет видеокарту на Windows
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

// detectGPUMacOS определяет видеокарту на macOS
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

// getNvidiaGPUUsage возвращает загруженность и температуру видеокарты NVIDIA
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

// getAMDGPUUsage возвращает загруженность и температуру видеокарты AMD
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

// getIntelGPUUsage возвращает загруженность и температуру видеокарты Intel
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

// parseFloat преобразует строку в число
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
