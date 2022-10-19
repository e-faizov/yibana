package handlers

import (
	"context"
	"fmt"
	"github.com/e-faizov/yibana/internal/interfaces"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/e-faizov/yibana/internal"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (s *storeTest) GetGauge(name string) (internal.Gauge, bool) {
	if name != s.gaugeName {
		return internal.Gauge(1), false
	}
	return s.gauge, true
}

func (s *storeTest) GetCounter(name string) (internal.Counter, bool) {
	if name != s.counterName {
		return internal.Counter(1), false
	}
	return s.counter, true
}

var gStore storeTest

func newRouter(h *MetricsHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update", h.PutJSONHandler)
	r.Post("/value", h.GetJSONHandler)

	r.Post("/update/{type}/{name}/{value}", h.PostHandler)
	r.Get("/value/{type}/{name}", h.GetHandler)

	return r
}

func testRequest(request *http.Request, store interfaces.Store) *http.Response {
	h := MetricsHandlers{
		Store: store,
	}
	request = request.WithContext(context.Background())
	w := httptest.NewRecorder()
	testHandle := newRouter(&h)
	testHandle.ServeHTTP(w, request)
	return w.Result()
}

func TestMetricsHandlers_UpdateWrongData(t *testing.T) {

	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "unknown type",
			request: "/update/unknown/test/1",
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			result := testRequest(request, &gStore)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}

}

func TestMetricsHandlers_Counters(t *testing.T) {
	type wantPost struct {
		statusCode int
		name       string
		value      internal.Counter
	}

	type wantGet struct {
		statusCode int
		name       string
		value      internal.Counter
	}

	type wantGetJSON struct {
		statusCode int
		name       string
		value      internal.Counter
	}

	key1 := "key1"
	key1first := 1
	key1second := 3242
	//key2 := "key2"

	tests := []struct {
		name    string
		request string
		method  string
		want    interface{}
	}{
		{
			name:    "set key1 first time",
			request: fmt.Sprintf("/update/counter/%s/%d", key1, key1first),
			method:  http.MethodPost,
			want: wantPost{
				statusCode: http.StatusOK,
				name:       key1,
				value:      internal.Counter(int64(key1first)),
			},
		},
		{
			name:    "get key1 first",
			request: fmt.Sprintf("/value/counter/%s", key1),
			method:  http.MethodGet,
			want: wantGet{
				statusCode: http.StatusOK,
				value:      internal.Counter(int64(key1first)),
			},
		},
		{
			name:    "get key1 first",
			request: fmt.Sprintf("/value/counter/%s", key1),
			method:  http.MethodGet,
			want: wantGet{
				statusCode: http.StatusOK,
				value:      internal.Counter(int64(key1first)),
			},
		},
		{
			name:    "set key1 by iter3 second time",
			request: fmt.Sprintf("/update/counter/%s/%d", key1, key1second),
			method:  http.MethodPost,
			want: wantPost{
				statusCode: http.StatusOK,
				name:       key1,
				value:      internal.Counter(int64(key1second)),
			},
		},
		{
			name:    "get key1 second",
			request: fmt.Sprintf("/value/counter/%s", key1),
			method:  http.MethodGet,
			want: wantGet{
				statusCode: http.StatusOK,
				value:      internal.Counter(int64(key1second)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			result := testRequest(request, &gStore)

			switch want := tt.want.(type) {
			case wantPost:
				assert.Equal(t, want.statusCode, result.StatusCode)
				err := result.Body.Close()
				require.NoError(t, err)
				assert.Equal(t, want.value, gStore.counter)
				assert.Equal(t, want.name, gStore.counterName)

			case wantGet:
				assert.Equal(t, want.statusCode, result.StatusCode)
				raw, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				res, err := strconv.ParseFloat(string(raw), 64)
				require.NoError(t, err)
				assert.Equal(t, want.value, internal.Counter(res))

			}
		})
	}
}

func TestMetricsHandlers_GetCounters(t *testing.T) {
	gStore.AddCounter("testCounter", 3534)

	type want struct {
		statusCode int
		value      string
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
			result := testRequest(request, &gStore)

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
	type want struct {
		statusCode int
		gaugeName  string
		gaugeValue internal.Gauge
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
			result := testRequest(request, &gStore)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.gaugeValue, gStore.gauge)
				assert.Equal(t, tt.want.gaugeName, gStore.gaugeName)
			}
		})
	}
}

func TestMetricsHandlers_GetGauges(t *testing.T) {
	gStore.SetGauge("testGauges", 3746.0)

	type want struct {
		statusCode int
		value      string
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
			result := testRequest(request, &gStore)

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
