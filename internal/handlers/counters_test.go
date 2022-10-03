package handlers

import (
	"github.com/e-faizov/yibana/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
func (s *storeTest) SetCounter(name string, val internal.Counter) error {
	s.counterName = name
	s.counter = val
	return nil
}

func TestMetricsHandlers_Counters(t *testing.T) {
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
			w := httptest.NewRecorder()
			h := http.HandlerFunc(h.Counters)
			h.ServeHTTP(w, request)
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
