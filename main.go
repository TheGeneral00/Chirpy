package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

type apiHandler struct {}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

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
    
    // Create and start server
    mux := http.NewServeMux()
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    //adding http.FileServer as handler with root /content
    mux.Handle("/", http.FileServer(http.Dir(".")))

    //adding a readiness endpoint
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
            w.Header().Set("Content-Type", "text/plain; charset=utf-8")
            w.WriteHeader(200)
            w.Write([]byte("OK"))
    })

    log.Printf("Starting server on port %v\n", server.Addr)
    err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Server error:", err)
    }
}
