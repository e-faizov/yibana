package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreFile     string        `env:"STORE_FILE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	Restore       bool          `env:"RESTORE"`
	KeyPath       string        `env:"CRYPTO_KEY"`
	TrustedSubnet string        `env:"TRUSTED_SUBNET"`
	GRPCPort      string        `env:"GRPC_PORT"`
}

type fileServerConfig struct {
	Address       *string        `json:"address,omitempty"`
	Restore       *bool          `json:"restore,omitempty"`
	StoreInterval *time.Duration `json:"store_interval,omitempty"`
	StoreFile     *string        `json:"store_file,omitempty"`
	DatabaseDsn   *string        `json:"database_dsn,omitempty"`
	CryptoKey     *string        `json:"crypto_key,omitempty"`
	TrustedSubnet *string        `json:"trusted_subnet,omitempty"`
	GRPCPort      *string        `json:"grpc_port,omitempty"`
}

var (
	serverCfg    ServerConfig
	srvCfgInited bool
)

func GetServerConfig() ServerConfig {
	if !srvCfgInited {
		readServerConfigFile()

		flag.StringVar(&serverCfg.Address, "a", "localhost:8080", "ADDRESS")
		flag.DurationVar(&serverCfg.StoreInterval, "i", time.Duration(300)*time.Second, "STORE_INTERVAL")
		flag.StringVar(&serverCfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE")
		flag.BoolVar(&serverCfg.Restore, "r", true, "RESTORE")
		flag.StringVar(&serverCfg.Key, "k", "", "KEY")
		flag.StringVar(&serverCfg.DatabaseDsn, "d", "", "KEY")
		flag.StringVar(&serverCfg.KeyPath, "crypto-key", "", "CRYPTO_KEY")
		flag.StringVar(&serverCfg.TrustedSubnet, "t", "", "TRUSTED_SUBNET")
		flag.StringVar(&serverCfg.GRPCPort, "gp", "", "GRPC_PORT")

		flag.Parse()
		if err := env.Parse(&serverCfg); err != nil {
			panic(err)
		}
		srvCfgInited = true
	}
	return serverCfg
}

func readServerConfigFile() error {
	configFile := getConfigFilePath()
	if len(configFile) != 0 {
		file, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}
		var fCfg fileServerConfig
		err = json.Unmarshal(file, &fCfg)
		if err != nil {
			return err
		}

		copyIfEnable(&serverCfg.Address, fCfg.Address)
		copyIfEnable(&serverCfg.Key, fCfg.CryptoKey)
		copyIfEnable(&serverCfg.StoreFile, fCfg.StoreFile)
		copyIfEnable(&serverCfg.DatabaseDsn, fCfg.DatabaseDsn)
		copyIfEnable(&serverCfg.TrustedSubnet, fCfg.TrustedSubnet)
		copyIfEnable(&serverCfg.GRPCPort, fCfg.GRPCPort)

		if fCfg.Restore != nil {
			serverCfg.Restore = *fCfg.Restore
		}

		if fCfg.StoreInterval != nil {
			serverCfg.StoreInterval = *fCfg.StoreInterval
		}

	}
	return nil
}
