package server

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/TheGeneral00/Chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
	JWTSecret      string
	PolkaKey       string
}

// New builds the HTTP server
func New(cfg *APIConfig, filepathRoot, port string) *http.Server {
	mux := http.NewServeMux()

	// File server with middlewares
	fsHandler := cfg.InputSanatizer(cfg.MiddlewareCreateUserEvent(cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
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

