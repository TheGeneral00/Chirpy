package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

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

    log.Printf("Starting server on port %v\n", server.Addr)
    err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Server error:", err)
    }
}
