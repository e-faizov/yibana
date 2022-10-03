package handlers

import (
	"fmt"
	"github.com/e-faizov/yibana/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
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

func (m *MetricsHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	if tp == "gauge" {
		res, err := m.Store.GetGauge(name)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%f", res)))
		return

	} else if tp == "counter" {
		res, err := m.Store.GetCounter(name)
		if err != nil {
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
