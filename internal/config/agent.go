// Package config - модуль для описания и  обработки конфигов риложений
package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// AgentConfig - конфиг для клиентского приложения
type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	Key            string        `env:"KEY"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	KeyPath        string        `env:"CRYPTO_KEY"`
}

var (
	agentCfg     AgentConfig
	agtCfgInited bool
)

type fileAgentConfig struct {
	Address        *string        `json:"address,omitempty"`
	ReportInterval *time.Duration `json:"report_interval,omitempty"`
	PollInterval   *time.Duration `json:"poll_interval,omitempty"`
	CryptoKey      *string        `json:"crypto_key,omitempty"`
}

func GetAgentConfig() AgentConfig {
	if !agtCfgInited {
		readAgentConfigFile()

		flag.StringVar(&(agentCfg.Address), "a", "localhost:8080", "ADDRESS")
		flag.DurationVar(&(agentCfg.ReportInterval), "r", time.Duration(10)*time.Second, "REPORT_INTERVAL")
		flag.DurationVar(&(agentCfg.PollInterval), "p", time.Duration(2)*time.Second, "POLL_INTERVAL")
		flag.StringVar(&(agentCfg.Key), "k", "", "KEY")
		flag.StringVar(&agentCfg.KeyPath, "crypto-key", "", "CRYPTO_KEY")

		flag.Parse()
		if err := env.Parse(&agentCfg); err != nil {
			panic(err)
		}

		agtCfgInited = true
	}
	return agentCfg
}

func readAgentConfigFile() error {
	configFile := getConfigFilePath()
	if len(configFile) != 0 {
		file, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}
		var fCfg fileAgentConfig
		err = json.Unmarshal(file, &fCfg)
		if err != nil {
			return err
		}

		if fCfg.Address != nil {
			agentCfg.Address = *fCfg.Address
		}

		if fCfg.ReportInterval != nil {
			agentCfg.ReportInterval = *fCfg.ReportInterval
		}

		if fCfg.PollInterval != nil {
			agentCfg.PollInterval = *fCfg.PollInterval
		}

		if fCfg.CryptoKey != nil {
			agentCfg.Key = *fCfg.CryptoKey
		}

	}
	return nil
}
