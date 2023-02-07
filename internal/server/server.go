package server

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/e-faizov/yibana/internal/encryption"
	"github.com/e-faizov/yibana/internal/handlers"
	"github.com/e-faizov/yibana/internal/interfaces"
	"github.com/e-faizov/yibana/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

// StartServer - функция запуска сервера
func StartServer(adr string, store interfaces.Store, key string, keyPath string) error {
	h := handlers.MetricsHandlers{
		Store: store,
		Key:   key,
	}

	r := chi.NewRouter()
	if len(keyPath) != 0 {
		privKey, err := encryption.ReadPrivKey(keyPath)
		if err != nil {
			return nil
		}
		r.Use(middlewares.DecryptFunc(privKey))
	}

	r.Use(middlewares.Compress)
	r.Use(middlewares.RequestLogger)

	r.Get("/", h.Info)
	r.Get("/ping", h.Ping)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.PutJSON)
		r.Post("/{type}/{name}/{value}", h.Post)
	})

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.PutsJSON)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetJSON)
		r.Get("/{type}/{name}", h.Get)
	})

	return http.ListenAndServe(adr, r)
}
