package internal

type Gauge float64
type Counter int64

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)
