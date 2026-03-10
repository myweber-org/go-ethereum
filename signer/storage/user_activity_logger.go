package main

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

var activityChannel = make(chan ActivityLog, 100)

func activityLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		activity := ActivityLog{
			UserID:    extractUserID(r),
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}

		select {
		case activityChannel <- activity:
		default:
			log.Println("Activity channel full, dropping log entry")
		}
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return "authenticated_user"
	}
	return "anonymous"
}

func processActivityLogs() {
	for activity := range activityChannel {
		log.Printf("Activity: %s %s by %s from %s at %v",
			activity.Method,
			activity.Endpoint,
			activity.UserID,
			activity.IPAddress,
			activity.Timestamp.Format(time.RFC3339))
	}
}

func main() {
	go processActivityLogs()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Data endpoint response"))
	})
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Info endpoint response"))
	})

	wrappedMux := activityLoggerMiddleware(mux)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatal(err)
	}
}