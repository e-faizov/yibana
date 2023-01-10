package handlers

import (
	"net/http"
	"strings"
	"testing"
)

func TestPuts(t *testing.T) {
	tStore := &storeTest{}

	bl := MetricsHandlers{
		Store: tStore,
	}

	testRouter := newRouter(&bl)
	t.Run("Empty body", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader("")
		req, err := http.NewRequest("POST", "/updates", body)
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
		req, err := http.NewRequest("POST", "/updates", body)
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
		req, err := http.NewRequest("POST", "/updates", body)
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
[
{
	"id": "test",
	"type": "gauge",
	"value": 100.1
},
{
	"id": "test2",
	"type": "gauge",
	"value": 435.3
}
]
`)
		tmpRt := newRouter(&MetricsHandlers{
			Store: tStore,
		})
		req, err := http.NewRequest("POST", "/updates", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(tmpRt, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})

	t.Run("wrong hash", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader(`
[{
	"id": "test",
	"type": "gauge",
	"value": 100.1,
	"hash": "13swq"
}]
`)
		tmpRt := newRouter(&MetricsHandlers{
			Store: tStore,
			Key:   "test",
		})
		req, err := http.NewRequest("POST", "/updates", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(tmpRt, req)

		if wr.Code != http.StatusBadRequest {
			t.Fatal("error, code not 400, code:", wr.Code)
		}
	})

	t.Run("wrong hash", func(t *testing.T) {
		tStore.Clear()
		body := strings.NewReader(`
[{
	"id": "test",
	"type": "gauge",
	"value": 100.1,
	"hash": "1a19f1c2dc35b7f40501c6c0dcea030f7a6b731bc92870fe4710dc91f520c604"
}]
`)
		tmpRt := newRouter(&MetricsHandlers{
			Store: tStore,
			Key:   "test",
		})
		req, err := http.NewRequest("POST", "/updates", body)
		if err != nil {
			t.Fatal(err)
		}

		wr := serveHTTP(tmpRt, req)

		if wr.Code != http.StatusOK {
			t.Fatal("error, code not 200, code:", wr.Code)
		}
	})
}
