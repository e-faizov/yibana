package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/e-faizov/yibana/internal/wg"
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

	ctxStop, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	var srv server.MetricsServer

	wg.Add()
	go func() {
		defer wg.Done()
		err = srv.StartServer(cfg, store)
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("can't start server")
		}
	}()

	<-ctxStop.Done()
	err = srv.Shutdown(ctxStop)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
