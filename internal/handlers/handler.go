package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/e-faizov/yibana/internal"
)

var (
	errNotFound    = errors.New("not found")
	errUnknownType = errors.New("unknown type")
	errWrongValue  = errors.New("wrong value")
	errSaveValue   = errors.New("error on save value")
)

func (m *MetricsHandlers) PutJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	//fmt.Println("save ", string(body))
	var data internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}

	err = m.putMetric(data)
	if err != nil {
		http.Error(w, errSaveValue.Error(), http.StatusBadRequest)
		return
	}
}

func (m *MetricsHandlers) GetJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	var data internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}
	ret, err := m.getValue(data.MType, data.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "error data format", http.StatusInternalServerError)
		return
	}

	w.Write(resData)
}

func (m *MetricsHandlers) putMetric(metric internal.Metric) error {
	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return errWrongValue
		}
		err := m.Store.SetGauge(metric.ID, *metric.Value)
		if err != nil {
			return errSaveValue
		}
	case "counter":
		if metric.Delta == nil {
			return errWrongValue
		}
		err := m.Store.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			return errSaveValue
		}
	default:
		return errUnknownType
	}
	return nil
}

func (m *MetricsHandlers) getValue(tp, key string) (internal.Metric, error) {
	ret := internal.Metric{
		ID: key,
	}

	if tp == "gauge" {
		res, ok := m.Store.GetGauge(key)
		if !ok {
			return internal.Metric{}, errNotFound
		}
		ret.SetGauge(res)
		return ret, nil

	} else if tp == "counter" {
		res, ok := m.Store.GetCounter(key)
		if !ok {
			return internal.Metric{}, errNotFound
		}
		ret.SetCounter(res)
		return ret, nil
	} else {
		return internal.Metric{}, errUnknownType
	}
}

func (m *MetricsHandlers) PostHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	data := internal.Metric{
		ID: name,
	}

	if tp == "gauge" {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, errWrongValue.Error(), http.StatusBadRequest)
			return
		}

		data.SetGauge(internal.Gauge(parsed))
	} else if tp == "counter" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, errWrongValue.Error(), http.StatusBadRequest)
			return
		}
		data.SetCounter(internal.Counter(parsed))
	} else {
		http.Error(w, errUnknownType.Error(), http.StatusNotImplemented)
		return
	}

	err := m.putMetric(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (m *MetricsHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")

	val, err := m.getValue(tp, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	switch val.MType {
	case "gauge":
		w.Write([]byte(fmt.Sprintf("%.3f", *val.Value)))
		return
	case "counter":
		w.Write([]byte(fmt.Sprintf("%d", *val.Delta)))
		return
	default:
		http.Error(w, "wrong path", http.StatusNotImplemented)
		return
	}
}
