package main

import (
	"time"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/config"
)

func main() {
	cfg := config.GetAgentConfig()

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)

	metrics := internal.Metrics{
		Key: cfg.Key,
	}
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
