package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

func NewActivityLog(userID, action, details string) *ActivityLog {
    return &ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }
}

func (al *ActivityLog) SaveToFile(filename string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    data, err := json.Marshal(al)
    if err != nil {
        return err
    }

    _, err = file.Write(append(data, '\n'))
    return err
}

func main() {
    logger := NewActivityLog("user123", "login", "User logged in from web browser")
    
    if err := logger.SaveToFile("activity.log"); err != nil {
        log.Fatal("Failed to save activity log:", err)
    }
    
    fmt.Println("Activity logged successfully")
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IP        string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)

		activity := ActivityLog{
			Timestamp: start,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
		}

		log.Printf("ACTIVITY: %s %s %s %s %d %v",
			activity.IP,
			activity.Method,
			activity.Path,
			activity.UserAgent,
			lw.statusCode,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
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
    Details   string    `json:"details,omitempty"`
}

type ActivityLogger struct {
    logFile string
}

func NewActivityLogger(logFile string) *ActivityLogger {
    return &ActivityLogger{logFile: logFile}
}

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now().UTC(),
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

func (al *ActivityLogger) ReadRecentActivities(limit int) ([]ActivityLog, error) {
    data, err := os.ReadFile(al.logFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []ActivityLog{}, nil
        }
        return nil, fmt.Errorf("failed to read log file: %w", err)
    }

    lines := []string{}
    currentLine := ""
    for _, b := range data {
        if b == '\n' {
            lines = append(lines, currentLine)
            currentLine = ""
        } else {
            currentLine += string(b)
        }
    }
    if currentLine != "" {
        lines = append(lines, currentLine)
    }

    var logs []ActivityLog
    start := len(lines) - limit
    if start < 0 {
        start = 0
    }

    for i := start; i < len(lines); i++ {
        if lines[i] == "" {
            continue
        }
        var logEntry ActivityLog
        if err := json.Unmarshal([]byte(lines[i]), &logEntry); err != nil {
            continue
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
    }

    err = logger.LogActivity("user123", "upload", "File 'report.pdf' uploaded successfully")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    activities, err := logger.ReadRecentActivities(5)
    if err != nil {
        fmt.Printf("Error reading activities: %v\n", err)
    }

    fmt.Println("Recent activities:")
    for _, activity := range activities {
        fmt.Printf("[%s] User %s: %s - %s\n",
            activity.Timestamp.Format("2006-01-02 15:04:05"),
            activity.UserID,
            activity.Action,
            activity.Details)
    }
}