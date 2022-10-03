package server

import (
	"fmt"
	"github.com/e-faizov/yibana/internal/handlers"
	"github.com/e-faizov/yibana/internal/storage"
	"net/http"
)

func StartServer(adr string, port int64) error {
	store := storage.NewStore()

	h := handlers.MetricsHandlers{
		Store: store,
	}

	http.HandleFunc("/update/", h.Handler)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", adr, port), nil)
}
