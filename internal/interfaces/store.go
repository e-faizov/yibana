package interfaces

import (
	"context"

	"github.com/e-faizov/yibana/internal"
)

// Store - интерфейс хранилища
type Store interface {
	// SetMetrics - сохранение списка метрик
	SetMetrics(ctx context.Context, metric []internal.Metric) error
	// GetMetric - получение метрики по имени и типу
	GetMetric(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error)
	// GetAll - получение всех метрик
	GetAll(ctx context.Context) ([]internal.Metric, error)
	// Ping - pong
	Ping() error
}
