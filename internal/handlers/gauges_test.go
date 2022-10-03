package handlers

import (
	"github.com/e-faizov/yibana/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			w := httptest.NewRecorder()
			h := http.HandlerFunc(h.Gauges)
			h.ServeHTTP(w, request)
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
