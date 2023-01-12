// Package config - модуль для описания и  обработки конфигов риложений
package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

// AgentConfig - конфиг для клиентского приложения
type AgentConfig struct {
	// Address - url сервиса сбора метрик
	Address string `env:"ADDRESS"`
	// ReportInterval - интервал отправки метрик
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	// PollInterval - интервал сбора метрик
	PollInterval time.Duration `env:"POLL_INTERVAL"`
	// Key - ключ хэширования
	Key string `env:"KEY"`
}

var (
	agentCfg     AgentConfig
	agtCfgInited bool
)

func GetAgentConfig() AgentConfig {
	if !agtCfgInited {
		flag.StringVar(&(agentCfg.Address), "a", "localhost:8080", "ADDRESS")
		flag.DurationVar(&(agentCfg.ReportInterval), "r", time.Duration(10)*time.Second, "REPORT_INTERVAL")
		flag.DurationVar(&(agentCfg.PollInterval), "p", time.Duration(2)*time.Second, "POLL_INTERVAL")
		flag.StringVar(&(agentCfg.Key), "k", "", "KEY")

		flag.Parse()
		if err := env.Parse(&agentCfg); err != nil {
			panic(err)
		}
		agtCfgInited = true
	}
	return agentCfg
}
