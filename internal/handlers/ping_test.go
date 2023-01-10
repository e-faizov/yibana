package handlers

import (
	"errors"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	tStore := &storeTest{}

	bl := MetricsHandlers{
		Store: tStore,
	}

	testRouter := newRouter(&bl)

	t.Run("error", func(t *testing.T) {
		tStore.Clear()
		tStore.ping = func() error {
			return errors.New("error")
		}
		req, err := http.NewRequest("GET", "/ping", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusInternalServerError {
			t.Fatal("error, code not 500, code:", wr.Code)
		}
	})

	t.Run("OK", func(t *testing.T) {
		tStore.Clear()
		tStore.ping = func() error {
			return nil
		}
		req, err := http.NewRequest("GET", "/ping", nil)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(testRouter, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})
}
