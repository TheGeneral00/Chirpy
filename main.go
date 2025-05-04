package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiHandler struct {}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type apiConfig struct {
        fileServerHits atomic.Int32
}

func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
                cfg.fileServerHits.Add(1)
                next.ServeHTTP(w, req)
        })
}

func main() {
    // Create log file
    logFile, err := os.OpenFile("server_log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
        os.Exit(1)
    }
    defer logFile.Close()

    // Send logs to both terminal and file
    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)
    
    log.Println("Logging system initialized")

    //initialize Config to track metrics 
    var apiCfg apiConfig
    
    // Create and start server
    mux := http.NewServeMux()
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    //adding http.FileServer as handler with root /content
    handler := http.FileServer(http.Dir("."))
    mux.Handle("/app/", apiCfg.middleWareMetricsInc(handler))

    //adding a readiness endpoint
    mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
            w.Header().Set("Content-Type", "text/plain; charset=utf-8")
            w.WriteHeader(200)
            w.Write([]byte("OK\n"))
    })

    // adding metrics endpoint
    mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, req *http.Request) {
            w.Header().Set("Content-Type", "text/plain; charset=utf-8")
            w.WriteHeader(200)
            hitCount := apiCfg.fileServerHits.Load()
            string := fmt.Sprintf("Hits: %v\n", hitCount)
            w.Write([]byte(string))
    })

    mux.HandleFunc("POST /reset", func(w http.ResponseWriter, req *http.Request) {
            apiCfg.fileServerHits.Store(0) 
            w.Write([]byte("Metrics have been resetted.\n"))
    })

    log.Printf("Starting server on port %v\n", server.Addr)
    err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Server error:", err)
    }
}
