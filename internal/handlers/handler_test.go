package handlers

import (
	"context"
	"errors"
	"github.com/e-faizov/yibana/internal"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type storeTest struct {
	gaugeName   string
	gauge       internal.Gauge
	counterName string
	counter     internal.Counter
}

func (s *storeTest) SetGauge(name string, val internal.Gauge) error {
	s.gaugeName = name
	s.gauge = val
	return nil
}
func (s *storeTest) AddCounter(name string, val internal.Counter) error {
	s.counterName = name
	s.counter = val
	return nil
}

func (s *storeTest) GetGauge(name string) (internal.Gauge, error) {
	if name != s.gaugeName {
		return internal.Gauge(1), errors.New("not found")
	}
	return s.gauge, nil
}

func (s *storeTest) GetCounter(name string) (internal.Counter, error) {
	if name != s.counterName {
		return internal.Counter(1), errors.New("not found")
	}
	return s.counter, nil
}

func TestMetricsHandlers_PostCounters(t *testing.T) {
	store := storeTest{}
	type want struct {
		statusCode int
		name       string
		value      internal.Counter
	}

	h := MetricsHandlers{
		Store: &store,
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "normal int value",
			request: "/update/counter/test/1",
			want: want{
				statusCode: http.StatusOK,
				name:       "test",
				value:      internal.Counter(int64(1)),
			},
		},
		{
			name:    "invalid value(string)",
			request: "/update/counter/test/string",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "invalid value(float)",
			request: "/update/counter/test/1.0",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "without value",
			request: "/update/counter/test",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request = request.WithContext(context.Background())
			w := httptest.NewRecorder()
			testHanle := chi.NewRouter()
			testHanle.Post("/update/{type}/{name}/{value}", h.PostHandler)
			testHanle.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.value, store.counter)
				assert.Equal(t, tt.want.name, store.counterName)
			}
		})
	}
}

func TestMetricsHandlers_GetCounters(t *testing.T) {
	store := storeTest{}

	store.AddCounter("testCounter", 3534)

	type want struct {
		statusCode int
		value      string
	}

	h := MetricsHandlers{
		Store: &store,
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "normal value",
			request: "/value/counter/testCounter",
			want: want{
				statusCode: http.StatusOK,
				value:      "3534",
			},
		},
		{
			name:    "unknown value",
			request: "/value/counter/notFound",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			request = request.WithContext(context.Background())
			w := httptest.NewRecorder()
			testHanle := chi.NewRouter()
			testHanle.Get("/value/{type}/{name}", h.GetHandler)
			testHanle.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			res, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.value, string(res))
			}
		})
	}
}

func TestMetricsHandlers_Gauges(t *testing.T) {
	store := storeTest{}
	type want struct {
		statusCode int
		gaugeName  string
		gaugeValue internal.Gauge
	}

	h := MetricsHandlers{
		Store: &store,
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "normal float value",
			request: "/update/gauge/test/1.0",
			want: want{
				statusCode: http.StatusOK,
				gaugeName:  "test",
				gaugeValue: internal.Gauge(1.0),
			},
		},
		{
			name:    "normal int value",
			request: "/update/gauge/test/1",
			want: want{
				statusCode: http.StatusOK,
				gaugeName:  "test",
				gaugeValue: internal.Gauge(1.0),
			},
		},
		{
			name:    "invalid value",
			request: "/update/gauge/test/string",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "without value",
			request: "/update/gauge/test",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request = request.WithContext(context.Background())
			w := httptest.NewRecorder()
			testHanle := chi.NewRouter()
			testHanle.Post("/update/{type}/{name}/{value}", h.PostHandler)
			testHanle.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.gaugeValue, store.gauge)
				assert.Equal(t, tt.want.gaugeName, store.gaugeName)
			}
		})
	}
}

func TestMetricsHandlers_GetGauges(t *testing.T) {
	store := storeTest{}

	store.SetGauge("testGauges", 3746.0)

	type want struct {
		statusCode int
		value      string
	}

	h := MetricsHandlers{
		Store: &store,
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "normal value",
			request: "/value/gauge/testGauges",
			want: want{
				statusCode: http.StatusOK,
				value:      "3746.000",
			},
		},
		{
			name:    "unknown value",
			request: "/value/gauge/notFound",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			request = request.WithContext(context.Background())
			w := httptest.NewRecorder()
			testHanle := chi.NewRouter()
			testHanle.Get("/value/{type}/{name}", h.GetHandler)
			testHanle.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			res, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.value, string(res))
			}
		})
	}
}
