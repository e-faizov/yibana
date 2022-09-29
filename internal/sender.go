package internal

import (
	"fmt"
	"net/http"
)

type Sender struct {
	adr  string
	port int64
}

func (s *Sender) SendMetrics(m Metrics) error {
	err := s.SendGauge(m.Alloc)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.BuckHashSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.Frees)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.GCCPUFraction)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.GCSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.HeapAlloc)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.HeapIdle)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.HeapInuse)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.HeapObjects)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.HeapReleased)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.HeapSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.LastGC)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.Lookups)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.MCacheInuse)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.MCacheSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.MSpanInuse)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.MSpanSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.Mallocs)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.NextGC)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.NumForcedGC)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.NumGC)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.OtherSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.PauseTotalNs)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.StackInuse)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.StackSys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.Sys)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.TotalAlloc)
	if err != nil {
		return err
	}

	err = s.SendCounter(m.PollCount)
	if err != nil {
		return err
	}

	err = s.SendGauge(m.RandomValue)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) SendGauge(gauge NamedGauge) error {
	url := fmt.Sprintf("http://%s:%d/update/gauge/%s/%f", s.adr, s.port, gauge.Name, gauge.Value)

	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) SendCounter(counter NamedCounter) error {
	url := fmt.Sprintf("http://%s:%d/update/counter/%s/%d", s.adr, s.port, counter.Name, counter.Value)

	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	return nil
}

func NewSender(adr string, port int64) Sender {
	return Sender{
		adr:  adr,
		port: port,
	}
}
