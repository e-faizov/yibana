package server

import (
	"fmt"
	"github.com/e-faizov/yibana/internal/handlers"
	"github.com/e-faizov/yibana/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func StartServer(adr string, port int64) error {
	store := storage.NewStore()

	h := handlers.MetricsHandlers{
		Store: store,
	}

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h.Handler)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", adr, port), r)
}
