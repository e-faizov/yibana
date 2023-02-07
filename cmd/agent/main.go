package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/config"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildVersion)
	fmt.Println("Build commit:", buildVersion)

	cfg := config.GetAgentConfig()

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)

	metrics := internal.Metrics{
		Key: cfg.Key,
	}
	if err := metrics.Update(); err != nil {
		log.Error().Err(err).Msg("error collection metrics")
	}

	sender, err := internal.NewSender(cfg.Address, cfg.KeyPath)
	if err != nil {
		panic(err)
	}

	go func() {
		for range pollTicker.C {
			if err := metrics.Update(); err != nil {
				log.Error().Err(err).Msg("error collection metrics")
			}
		}
	}()

	for range reportTicker.C {
		for {
			batch := metrics.Batch()
			if len(batch) != 0 {
				err := sender.SendMetrics(batch)
				if err != nil {
					log.Error().Err(err).Msg("can't send metrics")
				}
			}
		}
	}
}
