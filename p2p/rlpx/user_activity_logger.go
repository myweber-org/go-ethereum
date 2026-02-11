package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

type ActivityLogger struct {
    logFile string
}

func NewActivityLogger(logFile string) *ActivityLogger {
    return &ActivityLogger{logFile: logFile}
}

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile(al.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(logEntry); err != nil {
        return fmt.Errorf("failed to encode log entry: %w", err)
    }

    return nil
}

func (al *ActivityLogger) ReadLogs() ([]ActivityLog, error) {
    data, err := os.ReadFile(al.logFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []ActivityLog{}, nil
        }
        return nil, fmt.Errorf("failed to read log file: %w", err)
    }

    var logs []ActivityLog
    lines := bytes.Split(data, []byte("\n"))
    for _, line := range lines {
        if len(line) == 0 {
            continue
        }
        var logEntry ActivityLog
        if err := json.Unmarshal(line, &logEntry); err != nil {
            return nil, fmt.Errorf("failed to unmarshal log entry: %w", err)
        }
        logs = append(logs, logEntry)
    }

    return logs, nil
}

func main() {
    logger := NewActivityLogger("user_activity.json")
    
    err := logger.LogActivity("user123", "login", "User logged in from IP 192.168.1.100")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
        return
    }
    
    err = logger.LogActivity("user123", "upload", "File 'report.pdf' uploaded successfully")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
        return
    }
    
    logs, err := logger.ReadLogs()
    if err != nil {
        fmt.Printf("Error reading logs: %v\n", err)
        return
    }
    
    fmt.Printf("Total logged activities: %d\n", len(logs))
    for _, log := range logs {
        fmt.Printf("[%s] %s: %s - %s\n", 
            log.Timestamp.Format(time.RFC3339),
            log.UserID,
            log.Action,
            log.Details)
    }
}