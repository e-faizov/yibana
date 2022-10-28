package main

import (
	"github.com/e-faizov/yibana/internal/config"
	"github.com/e-faizov/yibana/internal/server"
	"github.com/e-faizov/yibana/internal/storage"
)

func main() {
	cfg := config.GetServerConfig()

	store, err := storage.NewStore(cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	if err != nil {
		panic(err)
	}

	err = server.StartServer(cfg.Address, store, cfg.Key)
	if err != nil {
		panic(err)
	}
}
