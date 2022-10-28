package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
)

func NewStore(storeInterval time.Duration, storeFile string, restore bool) (interfaces.Store, error) {
	var sync bool
	metrics := map[string]internal.Metric{}
	if restore {
		_, err := os.Stat(storeFile)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			data, err := os.ReadFile(storeFile)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &metrics)
			if err != nil {
				return nil, err
			}
		}
	}

	res := &storeImpl{
		metrics:   metrics,
		sync:      sync,
		storeFile: storeFile,
	}

	if storeInterval == 0 {
		sync = true
		return res, nil
	}

	dropper := time.NewTicker(storeInterval)
	go func() {
		for range dropper.C {
			err := res.Drop()
			if err != nil {
				fmt.Println("err drop", err.Error())
			}
		}
	}()

	return res, nil
}

type storeImpl struct {
	mtx     sync.RWMutex
	metrics map[string]internal.Metric

	storeInterval int
	storeFile     string
	sync          bool
}

func (s *storeImpl) SetMetric(metric internal.Metric) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if metric.MType == internal.GaugeType {
		s.metrics[metric.ID] = metric
	} else {
		old, ok := s.metrics[metric.ID]
		if ok && old.Delta != nil {
			*metric.Delta = *old.Delta + *metric.Delta
		}
		s.metrics[metric.ID] = metric
	}
	if s.sync {
		err := s.drop()
		if err != nil {
			fmt.Println("error drop", err.Error())
		}
	}
	return nil
}

func (s *storeImpl) GetMetric(metric internal.Metric) (internal.Metric, bool) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	res, ok := s.metrics[metric.ID]
	return res, ok
}

func (s *storeImpl) Drop() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.drop()
}

func (s *storeImpl) drop() error {
	data, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	return os.WriteFile(s.storeFile, data, 0644)
}

func (s *storeImpl) GetAll() []internal.Metric {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	res := make([]internal.Metric, 0, len(s.metrics))
	for _, v := range s.metrics {
		res = append(res, v)
	}
	return res
}
