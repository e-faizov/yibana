package handlers

import "net/http"

func (m *MetricsHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := m.Store.Ping()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
