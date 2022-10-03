package handlers

import (
	"github.com/e-faizov/yibana/internal"
	"net/http"
	"strconv"
	"strings"
)

func (m *MetricsHandlers) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 5 {
		http.Error(w, "wrong path", http.StatusNotFound)
		return
	}

	tp := strings.ToLower(paths[2])
	name := paths[3]

	if tp == "gauge" {
		val, err := strconv.ParseFloat(paths[4], 64)
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
		val, err := strconv.ParseInt(paths[4], 10, 64)
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
