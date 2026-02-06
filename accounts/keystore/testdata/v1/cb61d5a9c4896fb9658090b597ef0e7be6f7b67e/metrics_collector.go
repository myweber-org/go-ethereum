
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
	statusCodeCount map[int]int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		statusCodeCount: make(map[int]int),
	}
}

func (mc *MetricsCollector) RecordRequest(statusCode int, latency time.Duration) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.statusCodeCount[statusCode]++

	if statusCode >= 400 {
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

func (mc *MetricsCollector) GetStatusCodeDistribution() map[int]int {
	distribution := make(map[int]int)
	for code, count := range mc.statusCodeCount {
		distribution[code] = count
	}
	return distribution
}

func main() {
	collector := NewMetricsCollector()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metrics endpoint"))
		collector.RecordRequest(http.StatusOK, time.Since(start))
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		collector.RecordRequest(http.StatusInternalServerError, time.Since(start))
	})

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			log.Printf("Requests: %d, Avg Latency: %v, Error Rate: %.2f",
				collector.requestCount,
				collector.GetAverageLatency(),
				collector.GetErrorRate())
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}