package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/TheGeneral00/Chirpy/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Start logging
	logFile, err := server.InitializeLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to database
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	dbQueries := database.New(db)

	// Load secrets
	jwtSecret := os.Getenv("JWTSecret")
	polkaKey := os.Getenv("Polka_Key")

	// Build API config
	apiCfg := server.APIConfig{
		DBQueries: dbQueries,
		JWTSecret: jwtSecret,
		PolkaKey:  polkaKey,
	}

	const port = "8080"
	const filepathRoot = "./app"

	// Create server
	srv := server.New(&apiCfg, filepathRoot, port)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

