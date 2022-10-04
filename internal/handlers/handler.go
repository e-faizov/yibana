package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/e-faizov/yibana/internal"
)

func (m *MetricsHandlers) PostHandler(w http.ResponseWriter, r *http.Request) {

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

		err = m.Store.AddCounter(name, internal.Counter(val))
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

func (m *MetricsHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	if tp == "gauge" {
		res, ok := m.Store.GetGauge(name)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%.3f", res)))
		return

	} else if tp == "counter" {
		res, ok := m.Store.GetCounter(name)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%d", res)))
		return
	} else {
		http.Error(w, "wrong path", http.StatusNotImplemented)
		return
	}
}
