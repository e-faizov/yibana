package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/e-faizov/yibana/internal/server"
	"github.com/e-faizov/yibana/internal/storage"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("parse config error: %+v\n", err)
		return
	}

	store, err := storage.NewStore(cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	if err != nil {
		panic(err)
	}

	err = server.StartServer(cfg.Address, store)
	if err != nil {
		panic(err)
	}
}
