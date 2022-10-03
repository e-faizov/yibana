package interfaces

import "github.com/e-faizov/yibana/internal"

type Store interface {
	SetGauge(name string, val internal.Gauge) error
	SetCounter(name string, val internal.Counter) error

	GetGauge(name string) (internal.Gauge, error)
	GetCounter(name string) (internal.Counter, error)
}
