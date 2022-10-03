package storage

import (
	"errors"
	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
	"sync"
)

func NewStore() interfaces.Store {
	store := initStore()
	return &store
}

type storeImpl struct {
	gaugesMtx sync.RWMutex
	gauges    map[string]internal.Gauge

	countersMtx sync.RWMutex
	counters    map[string]internal.Counter
}

func initStore() storeImpl {
	return storeImpl{
		gauges:   map[string]internal.Gauge{},
		counters: map[string]internal.Counter{},
	}
}

var errNotFound = errors.New("not found")

func (s *storeImpl) SetGauge(name string, val internal.Gauge) error {
	s.gaugesMtx.Lock()
	defer s.gaugesMtx.Unlock()
	s.gauges[name] = val
	return nil
}
func (s *storeImpl) SetCounter(name string, val internal.Counter) error {
	s.countersMtx.Lock()
	defer s.countersMtx.Unlock()
	s.counters[name] = val
	return nil
}

func (s *storeImpl) GetGauge(name string) (internal.Gauge, error) {
	s.gaugesMtx.Lock()
	defer s.gaugesMtx.Unlock()
	v, ok := s.gauges[name]
	if !ok {
		return internal.Gauge(0), errNotFound
	}
	return v, nil
}

func (s *storeImpl) GetCounter(name string) (internal.Counter, error) {
	s.countersMtx.Lock()
	defer s.countersMtx.Unlock()
	v, ok := s.counters[name]
	if !ok {
		return internal.Counter(0), errNotFound
	}
	return v, nil
}
