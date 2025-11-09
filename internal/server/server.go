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

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	jwtSecret      string
	polkaKey       string
}

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

// Middleware to count server hits
func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// New builds the HTTP server
func New(cfg *apiConfig, filepathRoot, port string) *http.Server {
	mux := http.NewServeMux()

	// File server with middlewares
	fsHandler := cfg.MiddlewareCreateUserEvent(cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/", fsHandler)

	// Routes
	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/api/chirps", cfg.handlerAddChirp)
	mux.HandleFunc("/api/users", cfg.handlerAddUser)
	mux.HandleFunc("/api/login", cfg.handlerLogin)
	mux.HandleFunc("/api/revoke", cfg.handlerRevoke)
	mux.HandleFunc("/api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("/admin/reset", cfg.handlerReset)
	mux.HandleFunc("/admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/polka/webhooks", cfg.handlerPolkaWebhooks)
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

