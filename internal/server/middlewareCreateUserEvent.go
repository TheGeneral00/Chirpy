package server

import (
	"context"
	"net/http"

	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

// Middleware to create user event
func (cfg *APIConfig) MiddlewareCreateUserEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("X-User-ID")
		method := r.Method
		details := r.URL.Path

		userUUID, err := uuid.Parse(userId)

		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, "Invalid or missing X-User-ID", err)
			return
		}

		eventID, err := cfg.DBQueries.CreateUserEvent(r.Context(), database.CreateUserEventParams{
			UserID:        userUUID,
			Method:        method,
			MethodDetails: details,
		})
		if err != nil {
			cfg.Logger.Failure.Printf("Failed to store user event: %v", err)
		} else {
			cfg.Logger.Info.Printf("UserID: %v Method: %v URL: %v", userId, method, details)
		}

		ctx := context.WithValue(r.Context(), eventID, eventID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
