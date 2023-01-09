package handlers

import (
	"github.com/e-faizov/yibana/internal/interfaces"
)

// MetricsHandlers - структура для обработчиков сбора метрик
type MetricsHandlers struct {
	// Store - интерфейс хранилища
	Store interfaces.Store
	// Key - ключ для подсчета и проверки хэша
	Key string
}
