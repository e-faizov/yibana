package storage

import (
	"errors"
	"sync"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
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

var errNotFound = errors.New("not found")

func (s *storeImpl) SetGauge(name string, val internal.Gauge) error {
	s.gaugesMtx.Lock()
	defer s.gaugesMtx.Unlock()
	s.gauges[name] = val
	return nil
}
func (s *storeImpl) AddCounter(name string, val internal.Counter) error {
	s.countersMtx.Lock()
	defer s.countersMtx.Unlock()
	v := s.counters[name]
	s.counters[name] = v + val
	return nil
}

func (s *storeImpl) GetGauge(name string) (internal.Gauge, bool) {
	s.gaugesMtx.Lock()
	defer s.gaugesMtx.Unlock()
	v, ok := s.gauges[name]
	return v, ok
}

func (s *storeImpl) GetCounter(name string) (internal.Counter, bool) {
	s.countersMtx.Lock()
	defer s.countersMtx.Unlock()
	v, ok := s.counters[name]
	return v, ok
}
