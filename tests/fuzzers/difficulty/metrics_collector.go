package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	requestCount    int
	errorCount      int
	totalLatency    time.Duration
	statusCounts    map[int]int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		statusCounts: make(map[int]int),
	}
}

func (mc *MetricsCollector) RecordRequest(status int, latency time.Duration) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.statusCounts[status]++

	if status >= 400 {
		mc.errorCount++
	}
}

func (mc *MetricsCollector) GetAverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) GetErrorRate() float64 {
	if mc.requestCount == 0 {
		return 0.0
	}
	return float64(mc.errorCount) / float64(mc.requestCount)
}

func (mc *MetricsCollector) GetStatusDistribution() map[int]int {
	distribution := make(map[int]int)
	for k, v := range mc.statusCounts {
		distribution[k] = v
	}
	return distribution
}

func metricsMiddleware(next http.Handler, collector *MetricsCollector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		latency := time.Since(start)
		collector.RecordRequest(recorder.statusCode, latency)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := NewMetricsCollector()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := map[string]interface{}{
			"total_requests": collector.requestCount,
			"error_rate":     collector.GetErrorRate(),
			"avg_latency":    collector.GetAverageLatency().String(),
			"status_codes":   collector.GetStatusDistribution(),
		}
		
		// In production, use proper JSON marshaling
		w.Write([]byte(`{"status":"metrics endpoint"}`))
	})
	
	handler := metricsMiddleware(mux, collector)
	
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}