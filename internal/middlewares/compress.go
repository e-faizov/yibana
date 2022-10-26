package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

type gzipWriter struct {
	Rw      http.ResponseWriter
	Writer  *gzip.Writer
	written int64
}

var writers gzWriters

func (w *gzipWriter) WriteHeader(statusCode int) {
	w.Rw.WriteHeader(statusCode)
}
func (w gzipWriter) Header() http.Header {
	return w.Rw.Header()
}
func (w *gzipWriter) Write(b []byte) (int, error) {
	contentType := w.Header().Get("Content-Type")
	var err error

	if w.Writer == nil {
		gz := writers.Pop()
		if gz == nil {
			gz, err = gzip.NewWriterLevel(w.Rw, gzip.BestSpeed)
			if err != nil {
				return io.WriteString(w.Rw, err.Error())
			}
		} else {
			gz.Reset(w.Rw)
		}
		w.Writer = gz
	}

	if strings.Contains(contentType, "application/javascript") ||
		strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/css") ||
		strings.Contains(contentType, "text/html") ||
		strings.Contains(contentType, "text/plain") ||
		strings.Contains(contentType, "text/xml") {
		w.Header().Set("Content-Encoding", "gzip")
		return w.Writer.Write(b)
	}
	return w.Rw.Write(b)
}

func (w *gzipWriter) Close() {
	if w.Writer != nil {
		w.Header().Set("Content-Encoding", "gzip")
		w.Writer.Close()
		writers.Push(w.Writer)
		w.Writer = nil
	}
}

type gzWriters struct {
	mtx     sync.Mutex
	writers []*gzip.Writer
}

func (g *gzWriters) Pop() *gzip.Writer {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if len(g.writers) == 0 {
		return nil
	}
	ret := g.writers[0]
	g.writers = g.writers[0:]
	return ret
}

func (g *gzWriters) Push(writer *gzip.Writer) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.writers = append(g.writers, writer)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz := gzipWriter{Rw: w}
		defer gz.Close()
		next.ServeHTTP(&gz, r)
	})
}
