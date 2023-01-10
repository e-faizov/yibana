package handlers

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(m.Key)
		for _, metric := range data {
			if !checkHash(m.Key, metric) {
				log.Error().Err(err).
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

// checkHash - функция проверки хэша метрики
func checkHash(key string, metric internal.Metric) bool {
	if metric.MType == internal.GaugeType {
		if metric.Hash != internal.CalcGaugeHash(metric.ID, *metric.Value, key) {
			return false
		}
	} else {
		if metric.Hash != internal.CalcCounterHash(metric.ID, *metric.Delta, key) {
			return false
		}
	}
	return true
}
