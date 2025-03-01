package monitor

// Экспортируем переменные для использования в других пакетах
var (
	AlarmEnabled    bool                  // Флаг включения/выключения уведомлений
	AlarmThresholds = ThresholdSettings{} // Пороговые значения
	AlarmInterval   int                   // Интервал между уведомлениями (в минутах)
)
