package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/e-faizov/yibana/internal/server"
	"github.com/e-faizov/yibana/internal/storage"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Duration(300)*time.Second, "STORE_INTERVAL")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE")
	flag.BoolVar(&cfg.Restore, "r", true, "RESTORE")
}

func main() {
	flag.Parse()
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
