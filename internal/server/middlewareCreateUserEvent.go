package server

import (
	"context"
	"net/http"
	"time"

	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

type RequestMeta struct {
	RequestID 	uuid.UUID
	UserID		uuid.NullUUID
	Method		string
	Path 		string
	IP		string
	UserAgent	string
	StartedAt	time.Time
}

//only used for context 
type requestMetaKeyType struct{}

var requestMetaKey requestMetaKeyType

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
			return
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

		requestMeta := createRequestMeta(requestUUID, userUUID, r.Method, r.URL.Path, getRequestIP(r), r.UserAgent())
		ctx := context.WithValue(r.Context(), requestMetaKey, requestMeta)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createRequestMeta(requestID uuid.UUID, userID uuid.NullUUID, method , path , ip , userAgent string) RequestMeta {
	requestMeta := RequestMeta{
		RequestID: requestID,
		UserID: userID,
		Method: method,
		Path: path,
		IP: ip,
		UserAgent: userAgent,
		StartedAt: time.Now().UTC(),
	} 

	return requestMeta
}

//------ Helper functions for populating requestMeta struct ------

func getRequestIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
