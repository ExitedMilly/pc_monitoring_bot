package monitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

// IPInfo представляет информацию о IP-адресе
type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Org      string `json:"org"` // Провайдер
	Timezone string `json:"timezone"`
}

// ProcessTraffic представляет трафик процесса
type ProcessTraffic struct {
	Name       string
	DownloadMB float64
	UploadMB   float64
}

// TrafficStats представляет статистику трафика
type TrafficStats struct {
	DownloadMB float64
	UploadMB   float64
}

// GetIPInfo возвращает информацию о внешнем IP
func GetIPInfo() (IPInfo, error) {
	resp, err := http.Get("https://ipinfo.io")
	if err != nil {
		return IPInfo{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IPInfo{}, err
	}

	var info IPInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return IPInfo{}, err
	}

	return info, nil
}

// GetTopProcesses возвращает топ-3 процессов по использованию сети
func GetTopProcesses() ([]ProcessTraffic, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var traffic []ProcessTraffic
	for _, p := range processes {
		name, _ := p.Name()
		conns, _ := net.ConnectionsPid("all", p.Pid)

		var download, upload float64
		for _, conn := range conns {
			if conn.Status == "ESTABLISHED" {
				// Используем BytesRecv и BytesSent из IOCounters
				io, err := net.IOCounters(false)
				if err != nil {
					continue
				}
				download += float64(io[0].BytesRecv) / 1024 / 1024
				upload += float64(io[0].BytesSent) / 1024 / 1024
			}
		}

		traffic = append(traffic, ProcessTraffic{
			Name:       name,
			DownloadMB: download,
			UploadMB:   upload,
		})
	}

	// Сортируем по убыванию входящего трафика
	sort.Slice(traffic, func(i, j int) bool {
		return traffic[i].DownloadMB > traffic[j].DownloadMB
	})

	// Возвращаем топ-3
	if len(traffic) > 3 {
		return traffic[:3], nil
	}
	return traffic, nil
}

// GetNetworkSpeed возвращает текущую скорость сети
func GetNetworkSpeed() (float64, float64, error) {
	io1, err := net.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}

	time.Sleep(1 * time.Second) // Ждём 1 секунду

	io2, err := net.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}

	downloadSpeed := float64(io2[0].BytesRecv-io1[0].BytesRecv) / 1024 / 1024
	uploadSpeed := float64(io2[0].BytesSent-io1[0].BytesSent) / 1024 / 1024

	return downloadSpeed, uploadSpeed, nil
}

// GetTrafficLast5Min возвращает общий трафик за последние 5 минут
func GetTrafficLast5Min() (TrafficStats, error) {
	io, err := net.IOCounters(false)
	if err != nil {
		return TrafficStats{}, err
	}

	// TODO: Реализовать логику для сбора данных за 5 минут
	// Пока возвращаем текущие значения
	return TrafficStats{
		DownloadMB: float64(io[0].BytesRecv) / 1024 / 1024,
		UploadMB:   float64(io[0].BytesSent) / 1024 / 1024,
	}, nil
}

// GetNetworkUsage возвращает информацию о сети для команды /status
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
