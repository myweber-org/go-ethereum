package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type Activity struct {
    SessionID  string    `json:"session_id"`
    UserID     string    `json:"user_id"`
    Action     string    `json:"action"`
    Timestamp  time.Time `json:"timestamp"`
    IPAddress  string    `json:"ip_address"`
    UserAgent  string    `json:"user_agent"`
}

type ActivityLogger struct {
    activities []Activity
}

func NewActivityLogger() *ActivityLogger {
    return &ActivityLogger{
        activities: make([]Activity, 0),
    }
}

func (al *ActivityLogger) LogActivity(sessionID, userID, action, ip, agent string) {
    activity := Activity{
        SessionID: sessionID,
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        IPAddress: ip,
        UserAgent: agent,
    }
    al.activities = append(al.activities, activity)
    log.Printf("Activity logged: %s by user %s", action, userID)
}

func (al *ActivityLogger) GetUserActivities(userID string) []Activity {
    var userActivities []Activity
    for _, activity := range al.activities {
        if activity.UserID == userID {
            userActivities = append(userActivities, activity)
        }
    }
    return userActivities
}

func (al *ActivityLogger) GetRecentActivities(limit int) []Activity {
    if limit > len(al.activities) {
        limit = len(al.activities)
    }
    return al.activities[len(al.activities)-limit:]
}

func main() {
    logger := NewActivityLogger()

    http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        sessionID := r.Header.Get("X-Session-ID")
        userID := r.Header.Get("X-User-ID")
        action := r.URL.Query().Get("action")
        ip := r.RemoteAddr
        agent := r.UserAgent()

        if sessionID == "" || userID == "" || action == "" {
            http.Error(w, "Missing required parameters", http.StatusBadRequest)
            return
        }

        logger.LogActivity(sessionID, userID, action, ip, agent)
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintf(w, "Activity logged successfully")
    })

    http.HandleFunc("/activities", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        userID := r.URL.Query().Get("user_id")
        var activities []Activity

        if userID != "" {
            activities = logger.GetUserActivities(userID)
        } else {
            limit := 10
            activities = logger.GetRecentActivities(limit)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(activities)
    })

    log.Println("Starting activity logger server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}