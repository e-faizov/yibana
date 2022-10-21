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

func NewStore(StoreInterval time.Duration, StoreFile string, Restore bool) (interfaces.Store, error) {
	var sync bool
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

	res := &storeImpl{
		metrics:   metrics,
		sync:      sync,
		storeFile: StoreFile,
	}

	if StoreInterval == 0 {
		sync = true
	} else {
		dropper := time.NewTicker(StoreInterval)
		go func() {
			for range dropper.C {
				err := res.Drop()
				if err != nil {
					fmt.Println("err drop", err.Error())
				}
			}
		}()
	}

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
	/*var v interface{}
	if metric.MType == "gauge" {
		v = *metric.Value
	} else {
		v = *metric.Delta
	}
	fmt.Println("set metric", metric, "value", v)*/
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if metric.MType == "gauge" {
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
	//fmt.Println("read metric", metric)
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	res, ok := s.metrics[metric.ID]
	return res, ok
}

func (s *storeImpl) Drop() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.Drop()
}

func (s *storeImpl) drop() error {
	data, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	return os.WriteFile(s.storeFile, data, 0644)
}
