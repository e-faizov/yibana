package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Sender struct {
	adr  string
	port int64
}

func (s *Sender) SendMetric(m Metric) error {
	url := fmt.Sprintf("http://%s/update", s.adr)

	bd, err := json.Marshal(m)
	if err != nil {
		return ErrorHelper(fmt.Errorf("error json.Marshal %w", err))
	}

	return s.send(url, bd)
}

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

func NewSender(adr string) Sender {
	return Sender{
		adr: adr,
	}
}
