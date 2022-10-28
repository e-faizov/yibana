package main

import (
	"fmt"
	"github.com/e-faizov/yibana/internal/config"
	"github.com/e-faizov/yibana/internal/server"
	"github.com/e-faizov/yibana/internal/storage"
	"log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	cfg := config.GetServerConfig()

	store, err := storage.NewStore(cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	if err != nil {
		panic(err)
	}

	err = server.StartServer(cfg.Address, store, cfg.Key)
	if err != nil {
		fmt.Println(err)
	}
}
