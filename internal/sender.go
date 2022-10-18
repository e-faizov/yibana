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
	url := fmt.Sprintf("http://%s:%d/update", s.adr, s.port)

	bd, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(bd))
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	return nil
}

func NewSender(adr string, port int64) Sender {
	return Sender{
		adr:  adr,
		port: port,
	}
}
