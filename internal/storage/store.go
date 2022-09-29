package storage

import (
	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
	"sync"
)

func NewStore() interfaces.Store {
	return &storeImpl{
		gauges:   map[string]internal.Gauge{},
		counters: map[string]internal.Counter{},
	}
}

type storeImpl struct {
	gaugesMtx sync.RWMutex
	gauges    map[string]internal.Gauge

	countersMtx sync.RWMutex
	counters    map[string]internal.Counter
}

func (s *storeImpl) SetGauge(name string, val internal.Gauge) error {
	s.gaugesMtx.Lock()
	defer s.countersMtx.Unlock()
	s.gauges[name] = val
	return nil
}
func (s *storeImpl) SetCounter(name string, val internal.Counter) error {
	s.countersMtx.Lock()
	defer s.countersMtx.Unlock()
	s.counters[name] = val
	return nil
}
