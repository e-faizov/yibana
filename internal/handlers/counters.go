package handlers

import (
	"github.com/e-faizov/yibana/internal"
	"net/http"
	"strconv"
	"strings"
)

func (m *MetricsHandlers) Counters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	if len(paths) != 5 {
		http.Error(w, "wrong path", http.StatusNotFound)
		return
	}

	name := paths[3]
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
}
