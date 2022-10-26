package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
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

func (m *MetricsHandlers) PutJSON(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	//fmt.Println("call PutJSON, body", string(body))
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

func (m *MetricsHandlers) GetJSON(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}
	//fmt.Println("call GetJSON, body", string(body))
	var data internal.Metric
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "wrong body, not json", http.StatusBadRequest)
		return
	}
	ret, ok, err := m.getValue(data.MType, data.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, ret)
}

func (m *MetricsHandlers) putMetric(metric internal.Metric) error {
	switch metric.MType {
	case "gauge", "counter":
		err := m.Store.SetMetric(metric)
		if err != nil {
			return errSaveValue
		}
	default:
		return errUnknownType
	}
	return nil
}

func (m *MetricsHandlers) getValue(tp, key string) (internal.Metric, bool, error) {
	ret := internal.Metric{
		ID: key,
	}

	if tp == "gauge" || tp == "counter" {
		res, ok := m.Store.GetMetric(ret)
		return res, ok, nil
	} else {
		return internal.Metric{}, false, errUnknownType
	}
}

func (m *MetricsHandlers) Post(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	//fmt.Println("call Post, path", r.URL.Path)

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

func (m *MetricsHandlers) Get(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := chi.URLParam(r, "name")

	//fmt.Println("call Get, path", r.URL.Path)

	val, ok, err := m.getValue(tp, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
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
