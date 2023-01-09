package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
}

var (
	serverCfg    ServerConfig
	srvCfgInited bool
)

func GetServerConfig() ServerConfig {
	if !srvCfgInited {
		flag.StringVar(&serverCfg.Address, "a", "localhost:8080", "ADDRESS")
		flag.DurationVar(&serverCfg.StoreInterval, "i", time.Duration(300)*time.Second, "STORE_INTERVAL")
		flag.StringVar(&serverCfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE")
		flag.BoolVar(&serverCfg.Restore, "r", true, "RESTORE")
		flag.StringVar(&serverCfg.Key, "k", "", "KEY")
		flag.StringVar(&serverCfg.DatabaseDsn, "d", "postgresql://postgresUser:postgresPW@localhost:5455/postgresDB?sslmode=disable", "KEY")

		flag.Parse()
		if err := env.Parse(&serverCfg); err != nil {
			panic(err)
		}
		srvCfgInited = true
	}
	return serverCfg
}
