package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"net/http"
	"strings"

	"github.com/e-faizov/yibana/internal/encryption"
)

func DecryptFunc(privKey *rsa.PrivateKey) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if privKey != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				hash := sha256.New()
				unCrypt, err := encryption.DecryptOAEP(hash, rand.Reader, privKey, b, nil)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				r.ContentLength = int64(len(unCrypt))
				r.Body = io.NopCloser(bytes.NewReader(unCrypt))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Decrypt(next http.Handler) http.Handler {
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
