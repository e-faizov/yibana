package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
)

func NewStore(StoreInterval int, StoreFile string, Restore bool) (interfaces.Store, error) {
	var sync bool
	if StoreInterval == 0 {
		sync = true
	}
	metrics := map[string]internal.Metric{}
	if Restore {
		_, err := os.Stat(StoreFile)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			data, err := os.ReadFile(StoreFile)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &metrics)
			if err != nil {
				return nil, err
			}
		}
	}

	return &storeImpl{
		metrics:   metrics,
		sync:      sync,
		storeFile: StoreFile,
	}, nil
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
	s.metrics[metric.ID] = metric
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
	if s.sync {
		err := s.drop()
		if err != nil {
			fmt.Println("error drop", err.Error())
		}
	}
	return res, ok
}

func (s *storeImpl) drop() error {
	data, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	return os.WriteFile(s.storeFile, data, 0644)
}
