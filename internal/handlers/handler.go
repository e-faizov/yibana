package handlers

import (
	"github.com/e-faizov/yibana/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
)

func (m *MetricsHandlers) Handler(w http.ResponseWriter, r *http.Request) {

	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if tp == "gauge" {
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "wrong value", http.StatusBadRequest)
			return
		}

		err = m.Store.SetGauge(name, internal.Gauge(val))
		if err != nil {
			http.Error(w, "error on save value", http.StatusBadRequest)
			return
		}

		return
	} else if tp == "counter" {
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "wrong value", http.StatusBadRequest)
			return
		}

		err = m.Store.SetCounter(name, internal.Counter(val))
		if err != nil {
			http.Error(w, "error on save value", http.StatusBadRequest)
			return
		}
		return
	} else {
		http.Error(w, "wrong path", http.StatusNotImplemented)
		return
	}
}
