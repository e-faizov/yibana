package server

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e-faizov/yibana/internal/handlers"
	"github.com/e-faizov/yibana/internal/interfaces"
)

func StartServer(adr string, store interfaces.Store, key string) error {
	h := handlers.MetricsHandlers{
		Store: store,
		Key:   key,
	}

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Get("/", h.Info)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.PutJSON)
		r.Post("/{type}/{name}/{value}", h.Post)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetJSON)
		r.Get("/{type}/{name}", h.Get)
	})

	return http.ListenAndServe(adr, r)
}
