package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/e-faizov/yibana/internal/handlers"
	"github.com/e-faizov/yibana/internal/storage"
)

func StartServer(adr string, port int64) error {
	store := storage.NewStore()

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

	return http.ListenAndServe(fmt.Sprintf("%s:%d", adr, port), r)
}
