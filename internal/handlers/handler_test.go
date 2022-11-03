package handlers

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
)

type storeTest struct {
	metric *internal.Metric
}

func (s *storeTest) Ping() error {
	return nil
}

func (s *storeTest) SetMetrics(ctx context.Context, metric []internal.Metric) error {
	s.metric = &metric[0]
	return nil
}
func (s *storeTest) GetMetric(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
	if s.metric == nil {
		return internal.Metric{}, false, nil
	}
	return *s.metric, true, nil
}

func (s *storeTest) GetAll(ctx context.Context) ([]internal.Metric, error) {
	return []internal.Metric{}, nil
}

var gStore storeTest

func newRouter(h *MetricsHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update", h.PutJSON)
	r.Post("/value", h.GetJSON)

	r.Post("/update/{type}/{name}/{value}", h.Post)
	r.Get("/value/{type}/{name}", h.Get)

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

func TestMetricsHandlers_Errors(t *testing.T) {

	type want struct {
		statusCode int
	}
	var emptyStore storeTest

	gaugeTestData :=
		`{
"id": "testg",
"type": "gauge"
		}
`

	counterTestData :=
		`{
"id": "testc",
"type": "counter"
		}
`

	tests := []struct {
		name    string
		request string
		body    io.Reader
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
		{
			name:    "gauge not found",
			request: "/value/gauge/testg",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:    "gauge json not found",
			request: "/value/",
			body:    strings.NewReader(gaugeTestData),
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:    "counter not found",
			request: "/value/counter/testc",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:    "counter json not found",
			request: "/value/",
			body:    strings.NewReader(counterTestData),
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, tt.body)
			result := testRequest(request, &emptyStore)

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

	tests := []struct {
		name         string
		request      string
		method       string
		gaugesData   map[string]internal.Gauge
		countersData map[string]internal.Counter
		want         interface{}
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
				assert.Equal(t, want.value, *gStore.metric.Delta)
				assert.Equal(t, want.name, gStore.metric.ID)

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
	metr := internal.Metric{
		ID: "testCounter",
	}
	metr.SetCounter(3534)
	gStore.SetMetrics(context.Background(), []internal.Metric{metr})

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			result := testRequest(request, &gStore)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.gaugeValue, *gStore.metric.Value)
				assert.Equal(t, tt.want.gaugeName, gStore.metric.ID)
			}
		})
	}
}

func TestMetricsHandlers_GetGauges(t *testing.T) {
	metr := internal.Metric{
		ID: "testGauges",
	}
	metr.SetGauge(3746.0)
	gStore.SetMetrics(context.Background(), []internal.Metric{metr})

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
