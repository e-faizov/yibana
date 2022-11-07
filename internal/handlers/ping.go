package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (m *MetricsHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := m.Store.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Ping error store ping")
		http.Error(w, "", http.StatusInternalServerError)
	}
}
