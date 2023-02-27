package internal

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/e-faizov/yibana/internal/encryption"
	"github.com/e-faizov/yibana/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Sender interface {
	SendMetric(m Metric) error
	SendMetrics(m []Metric) error
}

// Sender - структура для отправки метрик на сервер
type HTTPSender struct {
	adr    string
	port   int64
	pubKey *rsa.PublicKey
}

// SendMetric - метод отправки одной метрики в формате json
func (s *HTTPSender) SendMetric(m Metric) error {
	url := fmt.Sprintf("http://%s/update", s.adr)

	bd, err := json.Marshal(m)
	if err != nil {
		return ErrorHelper(fmt.Errorf("error json.Marshal %w", err))
	}

	return s.send(url, bd)
}

// SendMetrics - метод отправки списка метрик в формате json
func (s *HTTPSender) SendMetrics(m []Metric) error {
	url := fmt.Sprintf("http://%s/updates", s.adr)

	bd, err := json.Marshal(m)
	if err != nil {
		return ErrorHelper(fmt.Errorf("error json.Marshal %w", err))
	}

	return s.send(url, bd)
}

func (s *HTTPSender) send(url string, data []byte) error {
	var err error
	if s.pubKey != nil {
		hash := sha256.New()
		data, err = encryption.EncryptOAEP(hash, rand.Reader, s.pubKey, data, nil)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	localIP, err := getLocalAddress()
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", localIP.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrorHelper(errors.New("status code not 200"))
	}
	return nil
}

// NewSender - функция создания нового объекта для отправки метрик
func NewHTTPSender(adr string, keyPath string) (Sender, error) {
	var rsaPubKey *rsa.PublicKey
	if len(keyPath) != 0 {
		var err error
		rsaPubKey, err = encryption.ReadPubKey(keyPath)
		if err != nil {
			return nil, err
		}
	}

	return &HTTPSender{
		adr:    adr,
		pubKey: rsaPubKey,
	}, nil
}

func NewGRPCSender(adr string) (Sender, error) {
	return &grpcSender{
		adr: adr,
	}, nil
}

type grpcSender struct {
	adr string
}

func (g *grpcSender) connection() (*grpc.ClientConn, error) {
	return grpc.Dial(g.adr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (g *grpcSender) SendMetric(m Metric) error {
	conn, err := g.connection()
	if err != nil {
		return err
	}
	defer conn.Close()
	c := proto.NewMetricsServiceClient(conn)

	resp, err := c.AddMetric(context.Background(), UnConvert(m))
	if err != nil {
		return err
	}

	if len(resp.Error) != 0 {
		return errors.New(resp.Error)
	}
	return nil
}
func (g *grpcSender) SendMetrics(m []Metric) error {
	conn, err := g.connection()
	if err != nil {
		return err
	}
	defer conn.Close()
	c := proto.NewMetricsServiceClient(conn)

	arr := make([]*proto.Metric, 0, len(m))

	for _, tmp := range m {
		arr = append(arr, UnConvert(tmp))
	}

	data := proto.Metrics{
		Data: arr,
	}

	resp, err := c.AddMetrics(context.Background(), &data)
	if err != nil {
		return err
	}

	if len(resp.Error) != 0 {
		return errors.New(resp.Error)
	}
	return nil
}
