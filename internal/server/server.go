package server

import (
	"github.com/e-faizov/yibana/internal/interfaces"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e-faizov/yibana/internal/handlers"
)

func StartServer(adr string, store interfaces.Store) error {
	h := handlers.MetricsHandlers{
		Store: store,
	}

	r := chi.NewRouter()
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
