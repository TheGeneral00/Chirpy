package main

import (
	"database/sql"
	"os"
	"log"
	"github.com/TheGeneral00/Chirpy/internal/database"
	"github.com/TheGeneral00/Chirpy/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//Start up logger to write user events and responses
	logger, eventFile, err := server.NewLog()
	if err != nil {
		panic(err)
	}
	defer eventFile.Close()

	//Modify std log module to write to stdout and errors.log 
	errFile, err := server.InitStdLogger()
	if err != nil {
		panic(err)
	}
	defer errFile.Close()

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
		PolkaKey: polkaKey,
		Logger: logger,
	}

	const port = "8080"
	const filepathRoot = "./app"

	// Create server
	srv := server.New(&apiCfg, filepathRoot, port)
	logger.Info.Printf("Serving files from %s:%s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

