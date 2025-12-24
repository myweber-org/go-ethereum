package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	requestCount    int
	totalLatency   time.Duration
	statusCodes    map[int]int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		statusCodes: make(map[int]int),
	}
}

func (m *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		m.requestCount++
		m.totalLatency += duration
		m.statusCodes[recorder.statusCode]++
	})
}

func (m *MetricsCollector) GetAverageLatency() time.Duration {
	if m.requestCount == 0 {
		return 0
	}
	return m.totalLatency / time.Duration(m.requestCount)
}

func (m *MetricsCollector) GetStatusCodeDistribution() map[int]int {
	result := make(map[int]int)
	for code, count := range m.statusCodes {
		result[code] = count
	}
	return result
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := NewMetricsCollector()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: collector.Middleware(mux),
	}
	
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			log.Printf("Requests: %d, Avg Latency: %v", collector.requestCount, collector.GetAverageLatency())
			for code, count := range collector.GetStatusCodeDistribution() {
				log.Printf("Status %d: %d", code, count)
			}
		}
	}()
	
	log.Fatal(server.ListenAndServe())
}