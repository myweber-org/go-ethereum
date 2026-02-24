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

func (mc *MetricsCollector) AverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) ErrorRate() float64 {
	if mc.requestCount == 0 {
		return 0.0
	}
	return float64(mc.errorCount) / float64(mc.requestCount)
}

func (mc *MetricsCollector) Reset() {
	mc.requestCount = 0
	mc.errorCount = 0
	mc.totalLatency = 0
	mc.statusCounts = make(map[int]int)
}

func main() {
	collector := NewMetricsCollector()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			latency := time.Since(start)
			collector.RecordRequest(http.StatusOK, latency)
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}