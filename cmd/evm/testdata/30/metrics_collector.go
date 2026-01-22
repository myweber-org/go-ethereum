package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	responseTimes []time.Duration
	totalRequests int
}

func (mc *MetricsCollector) RecordResponseTime(duration time.Duration) {
	mc.responseTimes = append(mc.responseTimes, duration)
	mc.totalRequests++
}

func (mc *MetricsCollector) AverageResponseTime() time.Duration {
	if mc.totalRequests == 0 {
		return 0
	}
	var total time.Duration
	for _, rt := range mc.responseTimes {
		total += rt
	}
	return total / time.Duration(mc.totalRequests)
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		mc.RecordResponseTime(duration)
	})
}

func main() {
	collector := &MetricsCollector{}
	mux := http.NewServeMux()
	mux.Handle("/api", collector.Middleware(http.HandlerFunc(apiHandler)))
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}