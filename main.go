package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
        dbQueries *database.Queries
        jwtSecret string
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
	}


	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerAddChirp)
        mux.HandleFunc("POST /api/users", apiCfg.handlerAddUser)
        mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
        mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
        mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	
        mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

