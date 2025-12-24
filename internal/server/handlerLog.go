package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

//------ Plain text log functiopns -----

type Logger struct {
	Info *log.Logger
	Warning *log.Logger 
	Failure *log.Logger
}

func ServerLog() (*Logger, *os.File, error) {
	//Create logs dir if not exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, err
	}
	//Build full path and open file 
	fullPath := "logs/server.log"
	logFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}

	multi := io.MultiWriter(os.Stdout, logFile)
	
	l := &Logger{
		Info: log.New(multi, "Info: ", log.LstdFlags|log.Lmicroseconds),
		Warning: log.New(multi, "Warning: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		Failure: log.New(multi, "Failure: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}

	l.Info.Println("Logger initialized.")

	return l, logFile, nil 
}


// Initialize the standard log package to write to errors.log and stdout
func ErrorLog() (*os.File, error) {
	// Check for existence of the logs dir is managed by function above
	fullPath := "logs/errors.log"
	logFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	//MultiWriter: file + stdout
	multi := io.MultiWriter(os.Stdout, logFile)

	//Configure the standard log package
	log.SetOutput(multi)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	//Optional: initial log entry
	//log.Println("Error logger initialized")

	return logFile, nil
}

func (cfg *APIConfig) LogSuccess(requestID uuid.UUID) {
	cfg.Logger.Info.Printf("Event %d terminated successfully.", requestID)
	cfg.DBQueries.SetStateSuccess(context.Background(), requestID)
}

func (cfg *APIConfig) LogFailure(requestID uuid.UUID, err error) {
	cfg.Logger.Failure.Printf("Event %d failed with error: %v", requestID, err)
	cfg.DBQueries.SetStateFailure(context.Background(), requestID) 
}

func (cfg *APIConfig) LogWarning(requestID uuid.UUID, message string) {
	cfg.Logger.Warning.Printf("Event %d Message: %s", requestID, message)
	cfg.DBQueries.SetStateSuccess(context.Background(), requestID)
}


//------ Json log functions ------

type LoggerJson struct {
	log	*log.Logger
}
type entryJson struct {
	Timestamp 	time.Time	`json:"timestamp"`
	Level		string		`json:"level"`
	RequestID	uuid.UUID	`json:"request_id"`
	UserID		uuid.NullUUID  	`json:"user_id"`
	Method		string		`json:"method"`
	Path		string		`json:"path"`
	State		string		`json:"state"`
	Message		string		`json:"message"`
}

func ServerLogJson() (*log.Logger, *os.File, error) {
	//Create files with .log.json extension
	logFile, err := os.OpenFile("logs/server.log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(logFile, "", 0) //no prefix, timestamp included in json 
	return logger, logFile, nil
}


/*------ main log wrapper functions ONLY USE THESE! ------
If needed adjust the under lying functions, but pay attention not to break normalization of log entries
*/

func (cfg *APIConfig) LogEvent(requestID uuid.UUID, userID uuid.NullUUID, method, path, state, msg, level string) {
	// Plain log with switch to differ by state
	switch state {
	
	case "Warning":
		cfg.Logger.Warning.Printf("level: %s request_id: %v user_id: %v method: %s path: %s state: %s message: %s",
			level, requestID, userID, method, path, state, msg)
	case "Failure":
		cfg.Logger.Failure.Printf("level: %s request_id: %v user_id: %v method: %s path: %s state: %s message: %s",
			level, requestID, userID, method, path, state, msg)
	default:
		cfg.Logger.Info.Printf("level: %s request_id: %v user_id: %v method: %s path: %s state: %s message: %s",
			level, requestID, userID, method, path, state, msg)
	}

	entry := entryJson{
		Timestamp: 	time.Now().UTC(),
		Level:		level,	
		RequestID: 	requestID,
		UserID: 	userID,
		Method: 	method,
		Path: 		path,
		State: 		state,
		Message: 	msg,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal json log entry for request %v \n Error: %v", requestID, err)
		return
	}

	cfg.LoggerJson.log.Println(string(jsonData))
}

func (cfg *APIConfig) logMissingRequestID(userID uuid.NullUUID, method, path string) {
	cfg.Logger.Failure.Printf("Failed to retrieve or create request id for:\n user_id: %v method: %s path: %s state: %s message: %s",
		userID, method, path, "Failure", "Event without request id.")

	entry := map[string]interface{}{
		"timestamp": 		time.Now().UTC(),
		"level": 		"Error",
		"request_id":		uuid.NullUUID{Valid: false},
		"user_id":		userID,
		"method":		method,
		"path": 		path,
		"state": 		"Failure",
		"message":		"Event without request id.",
	}
	entryData, err := json.Marshal(entry)
	
	if err != nil {
		log.Printf("Failed to marshal entry data: %v", err)
		return
	}

	cfg.LoggerJson.log.Println(string(entryData))
}
