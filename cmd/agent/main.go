package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"

	"github.com/e-faizov/yibana/internal"
)

type config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("parse config error: %+v\n", err)
		return
	}

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

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
