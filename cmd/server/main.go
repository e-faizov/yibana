package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/e-faizov/yibana/internal/config"
	"github.com/e-faizov/yibana/internal/interfaces"
	"github.com/e-faizov/yibana/internal/server"
	"github.com/e-faizov/yibana/internal/storage"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildVersion)
	fmt.Println("Build commit:", buildVersion)

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

	err = server.StartServer(cfg.Address, store, cfg.Key, cfg.KeyPath)
	if err != nil {
		log.Error().Err(err).Msg("can't start server")
	}
}
