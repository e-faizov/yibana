package handlers

import (
	"encoding/json"
	"github.com/e-faizov/yibana/internal"
	"io"
	"net/http"
)

func (m *MetricsHandlers) PutsJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}

	var data []internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}

	if len(m.Key) != 0 {
		for _, metric := range data {
			if !checkHash(m.Key, metric) {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		}
	}

	err = m.Store.SetMetrics(ctx, data)
	if err != nil {
		http.Error(w, errSaveValue.Error(), http.StatusBadRequest)
		return
	}
}

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
