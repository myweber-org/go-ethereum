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
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

func (mc *MetricsCollector) RecordRequest(latency time.Duration, isError bool) {
	mc.requestCount++
	mc.totalLatency += latency
	if isError {
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

func main() {
	collector := NewMetricsCollector()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			latency := time.Since(start)
			collector.RecordRequest(latency, false)
		}()

		avgLatency := collector.GetAverageLatency()
		errorRate := collector.GetErrorRate()

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"average_latency_ms":` + avgLatency.String() + `,"error_rate":` + string(errorRate) + `}`))
	})

	log.Println("Starting metrics server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}