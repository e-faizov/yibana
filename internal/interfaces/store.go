package interfaces

import (
	"context"

	"github.com/e-faizov/yibana/internal"
)

type Store interface {
	SetMetrics(ctx context.Context, metric []internal.Metric) error
	GetMetric(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error)
	GetAll(ctx context.Context) ([]internal.Metric, error)
	Ping() error
}
