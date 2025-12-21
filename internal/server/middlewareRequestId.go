package server

import (
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func middlewareCreateRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		requestID := uuid.New().String()
		r.Header.Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}
