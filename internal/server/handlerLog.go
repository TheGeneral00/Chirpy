package server

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Info *log.Logger
	Error *log.Logger
}

func NewLog() (*Logger, *os.File, error) {
	//Create logs dir if not exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, err
	}
	//Build full path and open file 
	fullPath := "logs/server.log"
	logFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, nil, err
	}

	multi := io.MultiWriter(os.Stdout, logFile)
	
	l := &Logger{
		Info: log.New(multi, "Info: ", log.LstdFlags|log.Lmicroseconds),
		Error: log.New(multi, "Error: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}

	l.Info.Println("Logger initialized.")

	return l, logFile, nil 
}


// Initialize the standard log package to write to errors.log and stdout
func InitStdLogger() (*os.File, error) {
	// Check for existence of the logs dir is managed by function above
	fullPath := "logs/errors.log"
	logFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	//MultiWriter: file + stdout
	multi := io.MultiWriter(os.Stdout, logFile)

	//Configure the standard log package
	log.SetOutput(multi)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	//Optional: initial log entry
	log.Println("Standard logger initialized")

	return logFile, nil
}


/*func initializeLogger(file string) (*os.File, error) {
        logFile, err := os.OpenFile(fmt.Sprintf("%s.txt", file), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
        if err != nil {
                return nil, fmt.Errorf("Failed to open log file: %v", err)
        }

        //Send logs to both terminal and file
        multiWriter := io.MultiWriter(os.Stdout, logFile)
        log.SetOutput(multiWriter)
        log.SetFlags(log.LstdFlags | log.Lmicroseconds)

        log.Printf("Loging system initialized")

        return logFile, nil
}*/
