package server

import (
	"fmt"
	"io"
	"log"
	"os"
)


func initializeLogger() (*os.File, error) {
        logFile, err := os.OpenFile("server_log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
        if err != nil {
                return nil, fmt.Errorf("Failed to open log file: %v", err)
        }

        //Send logs to both terminal and file
        multiWriter := io.MultiWriter(os.Stdout, logFile)
        log.SetOutput(multiWriter)
        log.SetFlags(log.LstdFlags | log.Lmicroseconds)

        log.Printf("Loging system initialized")

        return logFile, nil
}
