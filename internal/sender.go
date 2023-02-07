package internal

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/e-faizov/yibana/internal/encryption"
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
	var err error
	if s.pubKey != nil {
		hash := sha256.New()
		data, err = encryption.EncryptOAEP(hash, rand.Reader, s.pubKey, data, nil)
		if err != nil {
			return err
		}
	}
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
	var rsaPubKey *rsa.PublicKey
	if len(keyPath) != 0 {
		var err error
		rsaPubKey, err = encryption.ReadPubKey(keyPath)
		if err != nil {
			return Sender{}, err
		}

		fmt.Println(rsaPubKey)
	}

	return Sender{
		adr:    adr,
		pubKey: rsaPubKey,
	}, nil
}
