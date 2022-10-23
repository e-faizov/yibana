package interfaces

import "github.com/e-faizov/yibana/internal"

type Store interface {
	SetMetric(metric internal.Metric) error
	GetMetric(metric internal.Metric) (internal.Metric, bool)
	GetAll() []internal.Metric
}
