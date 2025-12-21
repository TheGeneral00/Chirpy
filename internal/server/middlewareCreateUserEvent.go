package server

import (
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

		var userUUID uuid.NullUUID

		if userId != "" {
			parsed, err := uuid.Parse(userId)
			if err != nil {
				cfg.Logger.Info.Printf("Invlid X-User-ID: %v\n", err)
				userUUID = uuid.NullUUID{Valid: false} //Store as Null 
			} else {
				userUUID = uuid.NullUUID{UUID: parsed, Valid: true}
			}
		} else {
			userUUID = uuid.NullUUID{Valid: false}
		}
	

		requestID := r.Header.Get("X-Request-ID")
		requestUUID, err := uuid.Parse(requestID) 
		if err != nil {
			cfg.Logger.Failure.Printf("Invalid or missing X-Request-ID: %v.\n", err)
			cfg.logMissingRequestID(userUUID, r.Method, r.URL.Path)
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to assign a requestID", nil)
		}

		_, err = cfg.DBQueries.CreateUserEvent(r.Context(), database.CreateUserEventParams{
			RequestID:	requestUUID,
			UserID:        	userUUID,
			Method:        	method,
			MethodDetails: 	details,
			EventSeq: 1,
		})
		if err != nil {
			cfg.Logger.Failure.Printf("Failed to store user event: %v", err)
		} else {
			cfg.Logger.Info.Printf("UserID: %v Method: %v URL: %v", userId, method, details)
		}

		next.ServeHTTP(w, r)
	})
}
