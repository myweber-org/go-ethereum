package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUUsage    float64
    MemoryAlloc uint64
    MemoryTotal uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now(),
        MemoryAlloc: m.Alloc,
        MemoryTotal: m.Sys,
        Goroutines:  runtime.NumGoroutine(),
        CPUUsage:    calculateCPUUsage(),
    }
}

func calculateCPUUsage() float64 {
    start := time.Now()
    runtime.Gosched()
    time.Sleep(50 * time.Millisecond)
    elapsed := time.Since(start)

    usage := (50.0 / (float64(elapsed.Milliseconds()) / 1000.0)) * 100
    if usage > 100 {
        return 100.0
    }
    return usage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] CPU: %.2f%% | Memory: %v/%v MB | Goroutines: %d\n",
        metrics.Timestamp.Format("15:04:05"),
        metrics.CPUUsage,
        metrics.MemoryAlloc/1024/1024,
        metrics.MemoryTotal/1024/1024,
        metrics.Goroutines,
    )
}

func main() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        metrics := collectMetrics()
        printMetrics(metrics)
    }
}package main

import (
	"log"
	"net/http"
	"time"
)

var (
	requestLatency = make(map[string]time.Duration)
	statusCodes    = make(map[int]int)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)
		requestLatency[r.URL.Path] = duration
		statusCodes[recorder.statusCode]++
		log.Printf("Request to %s took %v, status: %d", r.URL.Path, duration, recorder.statusCode)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"latency": "collected", "status_codes": "tracked"}`))
	})
	handler := metricsMiddleware(mux)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}