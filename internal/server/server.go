package server

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/TheGeneral00/Chirpy/internal/helpers"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
	JWTSecret      string
	PolkaKey       string
}

// Middleware to create user event
func (cfg *APIConfig) MiddlewareCreateUserEvent(next http.Handler) http.Handler {
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

		cfg.DBQueries.CreateUserEvent(r.Context(), database.CreateUserEventParams{
			UserID:        userUUID,
			Method:        method,
			MethodDetails: details,
		})
	})
}

// Middleware to count server hits
func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// New builds the HTTP server
func New(cfg *APIConfig, filepathRoot, port string) *http.Server {
	mux := http.NewServeMux()

	// File server with middlewares
	fsHandler := cfg.MiddlewareCreateUserEvent(cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/", fsHandler)

	// Routes
	mux.HandleFunc("/api/healthz", HandlerReadiness)
	mux.HandleFunc("/api/chirps", cfg.HandlerAddChirp)
	mux.HandleFunc("/api/users", cfg.HandlerAddUser)
	mux.HandleFunc("/api/login", cfg.HandlerLogin)
	mux.HandleFunc("/api/revoke", cfg.HandlerRevoke)
	mux.HandleFunc("/api/refresh", cfg.HandlerRefresh)
	mux.HandleFunc("/admin/reset", cfg.HandlerReset)
	mux.HandleFunc("/admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("/api/polka/webhooks", cfg.HandlerPolkaWebhooks)
	//mux.HandleFunc("Get /api/events", apiCfg.handlerGetAllEvents)
	//mux.HandleFunc("Get /api/events/latest", apiCfg.handlerGetLatestEvents)
	//mux.HandleFunc("Get /api/events/user/{userID}", apiCfg.handlerGetEventsByUser)
	//mux.HandleFunc("Get /api/events/method/{method}", apiCfg.handlerGetEventsByMethod)
	//mux.HandleFunc("Post /api/events/reset", apiCfg.handlerResetEvents)

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

// InitializeLogger moved from main for reuse
func InitializeLogger() (*os.File, error) {
	logFile, err := os.OpenFile("logs/server_log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	return logFile, nil
}

