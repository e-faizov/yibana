package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/e-faizov/yibana/internal"
)

var (
	// errNotFound - ошибка отсутствия значения
	errNotFound = errors.New("not found")
	// errUnknownType - ошибка неверного типа
	errUnknownType = errors.New("unknown type")
	// errWrongValue - ошибка неверного значения
	errWrongValue = errors.New("wrong value")
	// errSaveValue - ошибка при сохранении значения
	errSaveValue = errors.New("error on save value")
)

// PutJSON - обработчик сбора метрики в формате json
func (m *MetricsHandlers) PutJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("PutJSON error read body")
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}

	var data internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error().Err(err).Msg("PutJSON error unmarshal body")
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}

	if len(m.Key) != 0 {
		if !internal.CheckHash(m.Key, data) {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}

	err = m.putMetric(ctx, data)
	if err != nil {
		log.Error().Err(err).Msg("PutJSON error save metrics")
		http.Error(w, errSaveValue.Error(), http.StatusBadRequest)
		return
	}
}

// GetJSON - обработчик получения метрики в формате json
func (m *MetricsHandlers) GetJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("GetJSON error read body")
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}

	var data internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error().Err(err).Msg("GetJSON error unmarshal body")
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}
	ret, ok, err := m.getValue(ctx, data.MType, data.ID)
	if err != nil {
		log.Error().Err(err).Msg("GetJSON error get data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if len(m.Key) != 0 {
		if ret.MType == internal.GaugeType {
			ret.Hash = internal.CalcGaugeHash(ret.ID, *ret.Value, m.Key)
		} else {
			ret.Hash = internal.CalcCounterHash(ret.ID, *ret.Delta, m.Key)
		}
	}

	render.JSON(w, r, ret)
}

// putMetric - метод сохранения метрики в базу
func (m *MetricsHandlers) putMetric(ctx context.Context, metric internal.Metric) error {
	switch metric.MType {
	case internal.GaugeType, internal.CounterType:
		err := m.Store.SetMetrics(ctx, []internal.Metric{metric})
		if err != nil {
			return errSaveValue
		}
	default:
		return errUnknownType
	}
	return nil
}

// getValue - метод получения метрики из базы
func (m *MetricsHandlers) getValue(ctx context.Context, tp, key string) (internal.Metric, bool, error) {
	ret := internal.Metric{
		ID:    key,
		MType: tp,
	}

	if tp == internal.GaugeType || tp == internal.CounterType {
		res, ok, err := m.Store.GetMetric(ctx, ret)
		if err != nil {
			return res, ok, err
		}
		return res, ok, nil
	} else {
		return internal.Metric{}, false, errUnknownType
	}
}

// Post - обработчик сохранения метрики, через строку запроса
func (m *MetricsHandlers) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	data := internal.Metric{
		ID: name,
	}

	if tp == internal.GaugeType {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Error().Err(err).Msg("Post error parse float data type")
			http.Error(w, errWrongValue.Error(), http.StatusBadRequest)
			return
		}

		data.SetGauge(internal.Gauge(parsed))
	} else if tp == internal.CounterType {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Post error parse int data type")
			http.Error(w, errWrongValue.Error(), http.StatusBadRequest)
			return
		}
		data.SetCounter(internal.Counter(parsed))
	} else {
		http.Error(w, errUnknownType.Error(), http.StatusNotImplemented)
		return
	}

	err := m.putMetric(ctx, data)
	if err != nil {
		log.Error().Err(err).Msg("Post error save data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// Get - обработчик получения метрики через строку запроса
func (m *MetricsHandlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	if tp != internal.GaugeType && tp != internal.CounterType {
		http.Error(w, "wrong type", http.StatusNotImplemented)
		return
	}

	val, ok, err := m.getValue(ctx, tp, name)
	if err != nil {
		log.Error().Err(err).Msg("Get error read data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	switch val.MType {
	case internal.GaugeType:
		w.Write([]byte(fmt.Sprintf("%.3f", *val.Value)))
		return
	case internal.CounterType:
		w.Write([]byte(fmt.Sprintf("%d", *val.Delta)))
		return
	default:
		log.Error().Err(err).Msg("Get error wrong data type")
		http.Error(w, "wrong path", http.StatusNotImplemented)
		return
	}
}
