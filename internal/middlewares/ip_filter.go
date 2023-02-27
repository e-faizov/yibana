package middlewares

import (
	"net"
	"net/http"
)

func FilterIP(subNet *net.IPNet) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realIPStr := r.Header.Get("X-Real-IP")
			if len(realIPStr) == 0 {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			realIP := net.ParseIP(realIPStr)

			if !subNet.Contains(realIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
