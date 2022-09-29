package server

import (
	"fmt"
	"github.com/e-faizov/yibana/cmd/server/handlers"
	"github.com/e-faizov/yibana/cmd/server/storage"
	"net/http"
)

func StartServer(adr string, port int64) error {

	store := storage.NewStore()

	h := handlers.MetricsHandlers{
		Store: store,
	}

	http.HandleFunc("/update/gauge/", h.Gauges)
	http.HandleFunc("/update/counter/", h.Counters)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", adr, port), nil)
}
