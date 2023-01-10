package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Sender - структура для отправки метрик на сервер
type Sender struct {
	adr  string
	port int64
}

// SendMetric - метод отправки одной метрики в формате json
func (s *Sender) SendMetric(m Metric) error {
	url := fmt.Sprintf("http://%s/update", s.adr)

	bd, err := json.Marshal(m)
	if err != nil {
		return ErrorHelper(fmt.Errorf("error json.Marshal %w", err))
	}

	return s.send(url, bd)
}

// SendMetrics - метод отправки списка метрик в формате json
func (s *Sender) SendMetrics(m []Metric) error {
	url := fmt.Sprintf("http://%s/updates", s.adr)

	bd, err := json.Marshal(m)
	if err != nil {
		return ErrorHelper(fmt.Errorf("error json.Marshal %w", err))
	}

	return s.send(url, bd)
}

func (s *Sender) send(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return ErrorHelper(fmt.Errorf("error post data %w", err))
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrorHelper(errors.New("status code not 200"))
	}
	return nil
}

// NewSender - функция создания нового объекта для отправки метрик
func NewSender(adr string) Sender {
	return Sender{
		adr: adr,
	}
}
