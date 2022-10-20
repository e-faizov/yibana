package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/e-faizov/yibana/internal/server"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("parse config error: %+v\n", err)
		return
	}

	server.StartServer(cfg.Address)
}
