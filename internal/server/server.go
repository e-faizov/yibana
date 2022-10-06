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
	r.Post("/update/{type}/{name}/{value}", h.PostHandler)
	r.Get("/value/{type}/{name}", h.GetHandler)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", adr, port), r)
}
