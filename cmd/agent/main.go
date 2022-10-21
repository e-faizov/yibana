package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"

	"github.com/e-faizov/yibana/internal"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("parse config error: %+v\n", err)
		return
	}

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)

	metrics := internal.Metrics{}
	metrics.Update()

	sender := internal.NewSender(cfg.Address)

	go func() {
		for range pollTicker.C {
			metrics.Update()
		}
	}()

	for range reportTicker.C {
		for {
			next, ok := metrics.Front()
			if !ok {
				break
			}

			err := sender.SendMetric(next)
			if err == nil {
				metrics.Pop()
			}
		}
	}
}
