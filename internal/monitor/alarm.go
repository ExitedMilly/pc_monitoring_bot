package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ThresholdSettings содержит пороговые значения для уведомлений
type ThresholdSettings struct {
	CPUTemp      float64 `json:"cpu_temp"`      // Порог температуры CPU
	GPUTemp      float64 `json:"gpu_temp"`      // Порог температуры GPU
	CPUUsage     float64 `json:"cpu_usage"`     // Порог нагрузки CPU
	GPUUsage     float64 `json:"gpu_usage"`     // Порог нагрузки GPU
	MemoryUsage  float64 `json:"memory_usage"`  // Порог использования памяти
	NetworkUsage float64 `json:"network_usage"` // Порог использования сети
	DiskUsage    float64 `json:"disk_usage"`    // Порог загруженности дисков
}

var (
	thresholdsFile = "alarm_thresholds.json" // Файл для сохранения порогов
)

// LoadThresholds загружает пороговые значения из файла
func LoadThresholds() error {
	file, err := os.ReadFile(thresholdsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Файл с пороговыми значениями не найден. Создан новый.")
			return SaveThresholds() // Создаем файл, если он не существует
		}
		return err
	}
	return json.Unmarshal(file, &AlarmThresholds)
}

// SaveThresholds сохраняет пороговые значения в файл
func SaveThresholds() error {
	data, err := json.MarshalIndent(AlarmThresholds, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("Сохранение пороговых значений: %s", string(data))
	return os.WriteFile(thresholdsFile, data, 0644)
}

// StartAlarmMonitor запускает мониторинг для уведомлений
func StartAlarmMonitor(bot *tgbotapi.BotAPI, chatID int64) {
	var lastNotification time.Time // Время последнего уведомления

	for {
		if !AlarmEnabled {
			time.Sleep(10 * time.Second)
			continue
		}

		// Проверяем, прошло ли достаточно времени с последнего уведомления
		if time.Since(lastNotification) < time.Duration(AlarmInterval)*time.Minute {
			time.Sleep(10 * time.Second)
			continue
		}

		// Получаем текущие значения
		cpuTemp := GetCPUTempValue()       // float64
		gpuTemp := GetGPUTempValue()       // float64
		cpuUsage := GetCPUUsageValue()     // float64
		gpuUsage := GetGPUUsageValue()     // float64
		memUsage := GetMemoryUsageValue()  // float64
		netUsage := GetNetworkUsageValue() // float64
		diskUsage := GetDiskUsageValue()   // float64

		// Формируем уведомление
		var output strings.Builder
		output.WriteString("🚨 Внимание! Превышены пороговые значения:\n")

		// Проверяем каждое пороговое значение
		if AlarmThresholds.CPUTemp > 0 && cpuTemp > AlarmThresholds.CPUTemp {
			output.WriteString(fmt.Sprintf("🌡️ CPU: %.1f°C (порог: %.1f°C)\n", cpuTemp, AlarmThresholds.CPUTemp))
		}
		if AlarmThresholds.GPUTemp > 0 && gpuTemp > AlarmThresholds.GPUTemp {
			output.WriteString(fmt.Sprintf("🌡️ GPU: %.1f°C (порог: %.1f°C)\n", gpuTemp, AlarmThresholds.GPUTemp))
		}
		if AlarmThresholds.CPUUsage > 0 && cpuUsage > AlarmThresholds.CPUUsage {
			output.WriteString(fmt.Sprintf("⚙️ CPU: %.1f%% (порог: %.1f%%)\n", cpuUsage, AlarmThresholds.CPUUsage))
		}
		if AlarmThresholds.GPUUsage > 0 && gpuUsage > AlarmThresholds.GPUUsage {
			output.WriteString(fmt.Sprintf("🎮 GPU: %.1f%% (порог: %.1f%%)\n", gpuUsage, AlarmThresholds.GPUUsage))
		}
		if AlarmThresholds.MemoryUsage > 0 && memUsage > AlarmThresholds.MemoryUsage {
			output.WriteString(fmt.Sprintf("🧠 Память: %.1f%% (порог: %.1f%%)\n", memUsage, AlarmThresholds.MemoryUsage))
		}
		if AlarmThresholds.NetworkUsage > 0 && netUsage > AlarmThresholds.NetworkUsage {
			output.WriteString(fmt.Sprintf("🌐 Сеть: %.1f МБ/с (порог: %.1f МБ/с)\n", netUsage, AlarmThresholds.NetworkUsage))
		}
		if AlarmThresholds.DiskUsage > 0 && diskUsage > AlarmThresholds.DiskUsage {
			output.WriteString(fmt.Sprintf("💽 Диски: %.1f%% (порог: %.1f%%)\n", diskUsage, AlarmThresholds.DiskUsage))
		}

		// Если есть превышения и chatID не равен 0, отправляем уведомление
		if output.Len() > len("🚨 Внимание! Превышены пороговые значения:\n") && chatID != 0 {
			msg := tgbotapi.NewMessage(chatID, output.String())
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка при отправке уведомления: %v", err)
			} else {
				log.Println("Уведомление отправлено успешно.")
			}

			// Обновляем время последнего уведомления
			lastNotification = time.Now()
		}

		time.Sleep(10 * time.Second)
	}
}
