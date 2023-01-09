package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/e-faizov/yibana/internal"
)

func serveHTTP(handler *chi.Mux, req *http.Request) *httptest.ResponseRecorder {
	wr := httptest.NewRecorder()
	handler.ServeHTTP(wr, req)
	return wr
}

type storeTest struct {
	setMetrics func(ctx context.Context, metric []internal.Metric) error
	getMetric  func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error)
	getAll     func(ctx context.Context) ([]internal.Metric, error)
	ping       func() error
}

func (s *storeTest) Clear() {
	s.setMetrics = nil
	s.getMetric = nil
	s.getAll = nil
	s.ping = nil
}

func (s *storeTest) Ping() error {
	if s.ping != nil {
		return s.ping()
	}
	return nil
}

func (s *storeTest) SetMetrics(ctx context.Context, metric []internal.Metric) error {
	if s.setMetrics != nil {
		return s.setMetrics(ctx, metric)
	}
	return nil
}
func (s *storeTest) GetMetric(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
	if s.getMetric != nil {
		return s.getMetric(ctx, metric)
	}
	return internal.Metric{}, true, nil
}

func (s *storeTest) GetAll(ctx context.Context) ([]internal.Metric, error) {
	if s.getAll != nil {
		return s.getAll(ctx)
	}
	return []internal.Metric{}, nil
}

func newRouter(h *MetricsHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update", h.PutJSON)
	r.Post("/value", h.GetJSON)

	r.Get("/ping", h.Ping)

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.PutsJSON)
	})

	r.Post("/update/{type}/{name}/{value}", h.Post)
	r.Get("/value/{type}/{name}", h.Get)

	return r
}

func TestGetJSON(t *testing.T) {
	tStore := &storeTest{}

	bl := MetricsHandlers{
		Store: tStore,
	}

	testRouter := newRouter(&bl)

	t.Run("Empty body", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("")
		req, err := http.NewRequest("POST", "/value", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Body not json", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("{]")
		req, err := http.NewRequest("POST", "/value", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Body empty json", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("{}")
		req, err := http.NewRequest("POST", "/value", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		tStore.Clear()
		tStore.getMetric = func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
			return metric, false, nil
		}

		body := strings.NewReader(`
{
	"id": "test",
	"type": "gauge"
}
`)
		req, err := http.NewRequest("POST", "/value", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusNotFound {
			t.Fatal("error, code not 404, code:", wr.Code)
		}
	})

	t.Run("DB error", func(t *testing.T) {
		tStore.Clear()
		tStore.getMetric = func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
			return metric, true, errors.New("error")
		}

		body := strings.NewReader(`
{
	"id": "test",
	"type": "gauge"
}
`)
		req, err := http.NewRequest("POST", "/value", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})
}

func TestPutJSON(t *testing.T) {
	tStore := &storeTest{}

	bl := MetricsHandlers{
		Store: tStore,
	}

	testRouter := newRouter(&bl)

	t.Run("Empty body", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("")
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Body not json", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("{]")
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Body empty json", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("{}")
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Body empty json", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("{}")
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("OK without hash", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader(`
{
	"id": "test",
	"type": "gauge",
	"value": 100.1
}
`)
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})

	t.Run("wrong hash", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader(`
{
	"id": "test",
	"type": "gauge",
	"value": 100.1,
	"hash": "13swq"
}
`)
		bl.Key = "test"
		defer func() {
			bl.Key = ""
		}()
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("wrong hash", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader(`
{
	"id": "test",
	"type": "gauge",
	"value": 100.1,
	"hash": "1a19f1c2dc35b7f40501c6c0dcea030f7a6b731bc92870fe4710dc91f520c604"
}
`)
		bl.Key = "test"
		defer func() {
			bl.Key = ""
		}()
		req, err := http.NewRequest("POST", "/update", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})
}

func TestGet(t *testing.T) {
	tStore := &storeTest{}

	bl := MetricsHandlers{
		Store: tStore,
	}

	testRouter := newRouter(&bl)

	t.Run("Wrong type", func(t *testing.T) {
		tStore.Clear()
		tStore.getMetric = func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
			var val internal.Gauge = 10.1
			metric.Value = &val
			return metric, true, nil
		}
		req, err := http.NewRequest("GET", "/value/wrong/type", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusNotImplemented {
			t.Fatal("error, code not 501, code:", wr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		tStore.Clear()
		tStore.getMetric = func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
			return internal.Metric{}, false, errors.New("error")
		}
		req, err := http.NewRequest("GET", "/value/"+internal.GaugeType+"/type", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusInternalServerError {
			t.Fatal("error, code not 500, code:", wr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		tStore.Clear()
		tStore.getMetric = func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
			return metric, false, nil
		}
		req, err := http.NewRequest("GET", "/value/"+internal.GaugeType+"/type", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusNotFound {
			t.Fatal("error, code not 404, code:", wr.Code)
		}
	})

	t.Run("OK", func(t *testing.T) {
		tStore.Clear()
		tStore.getMetric = func(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
			var val internal.Gauge = 10.1
			metric.Value = &val
			return metric, true, nil
		}
		req, err := http.NewRequest("GET", "/value/"+internal.GaugeType+"/type", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}

		body, err := io.ReadAll(wr.Body)
		if err != nil {
			t.Error("response error read body", err)
			return
		}

		if string(body) != "10.100" {
			t.Fatal("error, wrong value, body:", string(body))
		}
	})
}

func TestPut(t *testing.T) {
	tStore := &storeTest{}

	bl := MetricsHandlers{
		Store: tStore,
	}

	testRouter := newRouter(&bl)

	t.Run("Wrong type", func(t *testing.T) {
		tStore.Clear()
		req, err := http.NewRequest("POST", "/update/wrong/type/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusNotImplemented {
			t.Fatal("error, code not 501, code:", wr.Code)
		}
	})

	t.Run("Gauge wrong type", func(t *testing.T) {
		tStore.Clear()
		req, err := http.NewRequest("POST", "/update/"+internal.GaugeType+"/test/value", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Gauge DB error", func(t *testing.T) {
		tStore.Clear()
		tStore.setMetrics = func(ctx context.Context, metric []internal.Metric) error {
			return errors.New("error")
		}
		req, err := http.NewRequest("POST", "/update/"+internal.GaugeType+"/test/123.1", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Gauge OK float", func(t *testing.T) {
		tStore.Clear()
		tStore.setMetrics = func(ctx context.Context, metric []internal.Metric) error {
			return nil
		}
		req, err := http.NewRequest("POST", "/update/"+internal.GaugeType+"/test/123.1", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})

	t.Run("Gauge OK int", func(t *testing.T) {
		tStore.Clear()
		tStore.setMetrics = func(ctx context.Context, metric []internal.Metric) error {
			return nil
		}
		req, err := http.NewRequest("POST", "/update/"+internal.GaugeType+"/test/123", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})

	t.Run("Counter wrong types", func(t *testing.T) {
		tStore.Clear()
		req, err := http.NewRequest("POST", "/update/"+internal.CounterType+"/test/value", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("Counter OK int", func(t *testing.T) {
		tStore.Clear()
		req, err := http.NewRequest("POST", "/update/"+internal.CounterType+"/test/10", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})

	t.Run("Counter DB error", func(t *testing.T) {
		tStore.Clear()
		tStore.setMetrics = func(ctx context.Context, metric []internal.Metric) error {
			return errors.New("error")
		}
		req, err := http.NewRequest("POST", "/update/"+internal.CounterType+"/test/10", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

}
