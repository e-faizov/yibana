package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/e-faizov/yibana/internal"
)

// PutsJSON - обработчик для сохранения списка метрик
func (m *MetricsHandlers) PutsJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("PutsJSON error read body")
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}

	var data []internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error().Err(err).Msg("PutsJSON error unmarshal body")
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}

	if len(m.Key) != 0 {
		for _, metric := range data {
			if !internal.CheckHash(m.Key, metric) {
				log.Error().
					Str("hash", metric.Hash).
					Str("id", metric.ID).
					Msg("PutsJSON error wrong hash")
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		}
	}

	err = m.Store.SetMetrics(ctx, data)
	if err != nil {
		log.Error().Err(err).Msg("PutsJSON error save data")
		http.Error(w, errSaveValue.Error(), http.StatusBadRequest)
		return
	}
}
