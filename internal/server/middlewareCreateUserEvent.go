package server

import (
	"net/http"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

// Middleware to create user event
func (cfg *apiConfig) MiddlewareCreateUserEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("X-User-ID")
		method := r.Method
		details := r.URL.Path

		next.ServeHTTP(w, r)

		userUUID, err := uuid.Parse(userId)
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, "Invalid or missing X-User-ID", err)
			return
		}

		cfg.dbQueries.CreateUserEvent(r.Context(), database.CreateUserEventParams{
			UserID:        userUUID,
			Method:        method,
			MethodDetails: details,
		})
	})
}
