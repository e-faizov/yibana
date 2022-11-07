package handlers

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func (m *MetricsHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := m.Store.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Ping error store ping")
		http.Error(w, "", http.StatusInternalServerError)
	}
}
