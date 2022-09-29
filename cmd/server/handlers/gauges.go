package handlers

import (
	"github.com/e-faizov/yibana/internal"
	"net/http"
	"strconv"
	"strings"
)

func (m *MetricsHandlers) Gauges(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 4 {
		http.Error(w, "wrong path", http.StatusBadRequest)
		return
	}

	name := paths[2]
	val, err := strconv.ParseFloat(paths[3], 64)
	if err != nil {
		http.Error(w, "wrong value", http.StatusBadRequest)
		return
	}

	err = m.Store.SetGauge(name, internal.Gauge(val))
	if err != nil {
		http.Error(w, "error on save value", http.StatusBadRequest)
		return
	}
}
