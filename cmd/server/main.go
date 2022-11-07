package main

import (
	"github.com/rs/zerolog/log"

	"github.com/e-faizov/yibana/internal/config"
	"github.com/e-faizov/yibana/internal/interfaces"
	"github.com/e-faizov/yibana/internal/server"
	"github.com/e-faizov/yibana/internal/storage"
)

func main() {
	cfg := config.GetServerConfig()

	var store interfaces.Store
	var err error
	if len(cfg.DatabaseDsn) != 0 {
		store, err = storage.NewPgStore(cfg.DatabaseDsn)
	} else {
		store, err = storage.NewMemStore(cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	}
	if err != nil {
		panic(err)
	}

	err = server.StartServer(cfg.Address, store, cfg.Key)
	if err != nil {
		log.Error().Err(err).Msg("can't start server")
	}
}
