package internal

import (
	"errors"
	"fmt"
	"net/http"
)

type Sender struct {
	adr  string
	port int64
}

func (s *Sender) SendMetric(m Metric) error {
	return s.sendData(m.GetType(), m.Name, m.ToString())
}

func (s *Sender) sendData(tp, name, val string) error {
	url := fmt.Sprintf("http://%s:%d/update/%s/%s/%s", s.adr, s.port, tp, name, val)
	return s.send(url)
}

func (s *Sender) send(url string) error {
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
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
