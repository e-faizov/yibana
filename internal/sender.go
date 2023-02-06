package internal

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh"
)

// Sender - структура для отправки метрик на сервер
type Sender struct {
	adr    string
	port   int64
	pubKey *rsa.PublicKey
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
func NewSender(adr string, keyPath string) (Sender, error) {
	bytes, err := os.ReadFile(keyPath)
	if err != nil {
		return Sender{}, err
	}
	pubKey, err := ssh.ParsePublicKey(bytes)
	if err != nil {
		return Sender{}, err
	}

	fmt.Println(pubKey.Type())

	return Sender{
		adr: adr,
	}, nil
}
