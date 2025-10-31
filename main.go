package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
        dbQueries *database.Queries
        jwtSecret string
        polkaKey string
}

func (cfg *apiConfig) middlewareCreateUserEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		userId := r.Header.Get("X-User-ID")
		method := r.Method
		details := r.URL.Path

		next.ServeHTTP(w, r)
			
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid or missing X-User-ID", err)
			return
		}

		cfg.dbQueries.CreateUserEvent(r.Context(), database.CreateUserEventParams{
			UserID: userUUID,
			Method: method,
			MethodDetails: details,
		})
	})
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//keep count of server hits 
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
        //Starting up log 
        logFile, err := initializeLogger()
        if err != nil {
                fmt.Fprintf(os.Stderr, "%v\n", err)
                os.Exit(1)
        }
        defer logFile.Close()

        //calling godotenv to auto import env 
        err = godotenv.Load()
        if err != nil{
                log.Fatalf("Error loading .env file: %v", err)
        }

        //Set connection to database
        dbURL := os.Getenv("DB_URL")
        db, err := sql.Open("postgres", dbURL)
        dbQueries := database.New(db)

        const filepathRoot = "./app"
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
                dbQueries: dbQueries,
                jwtSecret: os.Getenv("JWTSecret"),
                polkaKey: os.Getenv("Polka_Key"),
	}

	


	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareCreateUserEvent(apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerAddChirp)
        mux.HandleFunc("POST /api/users", apiCfg.handlerAddUser)
        mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
        mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
        mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
        mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
        mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
        mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
        mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
        mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUserCredentials)
        mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)
	mux.HandleFunc("Get /api/events", apiCfg.handlerGetAllEvents)
	mux.HandleFunc("Get /api/events/latest", apiCfg.handlerGetLatestEvents)
	mux.HandleFunc("Get /api/events/user/{userID}", apiCfg.handlerGetEventsByUser)
	mux.HandleFunc("Get /api/events/method/{method}", apiCfg.handlerGetEventsByMethod)
	mux.HandleFunc("Post /api/events/reset", apiCfg.handlerResetEvents)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
