package api

import (
	"crypto/subtle"
	"log"
	"net/http"
	"strings"
)

func (h *Handlers) authMiddleware(next http.Handler) http.Handler {
	if h.apiKey == "" {
		log.Println("WARNING: RUPTURA_API_KEY is not set — all API requests are unauthenticated")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") ||
			subtle.ConstantTimeCompare([]byte(auth[7:]), []byte(h.apiKey)) != 1 {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
