package internal

import (
	"fmt"
	"net/http"
)

type Sender struct {
	adr  string
	port int64
}

func (s *Sender) SendMetric(m Metric) error {
	s.Send(m.GetType(), m.Name, m.ToString())
	return nil
}

func (s *Sender) Send(tp, name, val string) error {
	url := fmt.Sprintf("http://%s:%d/update/%s/%s/%s", s.adr, s.port, tp, name, val)
	return s.send(url)
}

func (s *Sender) send(url string) error {
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

func NewSender(adr string, port int64) Sender {
	return Sender{
		adr:  adr,
		port: port,
	}
}
