package functions

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shirou/gopsutil/process"
)

// processCache хранит процессы для каждого чата (ключ - chatID)
var (
	processCache = make(map[int64][]*process.Process)
	cacheMutex   sync.Mutex
)

// GetProcesses возвращает процессы для чата (из кэша или загружает)
func GetProcesses(chatID int64) ([]*process.Process, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if procs, ok := processCache[chatID]; ok {
		return procs, nil
	}

	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	processCache[chatID] = procs
	return procs, nil
}

// GetProcessesPage возвращает процессы для страницы и общее количество страниц
func GetProcessesPage(chatID int64, page int) ([]*process.Process, int, error) {
	procs, err := GetProcesses(chatID)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (len(procs) + 39) / 40 // Округление вверх
	start := page * 40
	end := start + 40

	if start >= len(procs) {
		return []*process.Process{}, totalPages, nil
	}

	if end > len(procs) {
		end = len(procs)
	}

	return procs[start:end], totalPages, nil
}

// FormatProcessesMessage формирует сообщение со списком процессов и клавиатуру для пагинации
func FormatProcessesMessage(procs []*process.Process, page, totalPages int) (string, tgbotapi.InlineKeyboardMarkup) {
	var sb strings.Builder

	// Формируем список процессов
	for i, p := range procs {
		name, err := p.Name()
		if err != nil {
			name = "unknown" // Если не удалось получить имя процесса
		}
		pid := p.Pid // Получаем PID процесса
		sb.WriteString(fmt.Sprintf(
			"%d. %s [PID: %d]\n",
			(page*40)+i+1, // Порядковый номер с учетом страницы
			name,          // Имя процесса
			pid,           // PID процесса
		))
	}

	// Добавляем информацию о текущей странице и общем количестве страниц
	sb.WriteString(fmt.Sprintf("\nСтраница %d/%d", page+1, totalPages))

	// Создаем клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	// Если есть следующая страница, добавляем кнопку "Вывести ещё"
	if page+1 < totalPages {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вывести ещё", "showproc_page_"+strconv.Itoa(page+1)),
		))
	}

	// Возвращаем сформированное сообщение и клавиатуру
	return sb.String(), keyboard
}

// HandleShowProcCommand обрабатывает команду /showproc
func HandleShowProcCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID

	// Получаем первую страницу
	procs, totalPages, err := GetProcessesPage(chatID, 0)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при получении процессов")
		bot.Send(msg)
		return
	}

	// Формируем сообщение
	message, keyboard := FormatProcessesMessage(procs, 0, totalPages)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
