package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"

	"github.com/e-faizov/yibana/internal"
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Duration(10)*time.Second, "REPORT_INTERVAL")
	flag.DurationVar(&cfg.PollInterval, "p", time.Duration(2)*time.Second, "POLL_INTERVAL")
}

func main() {
	flag.Parse()
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
