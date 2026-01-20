package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     string  `json:"timestamp"`
    CPUUsage      float64 `json:"cpu_usage"`
    MemoryUsage   uint64  `json:"memory_usage"`
    GoroutineCount int    `json:"goroutine_count"`
    AllocMemory   uint64 `json:"alloc_memory"`
    TotalMemory   uint64 `json:"total_memory"`
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now().UTC().Format(time.RFC3339),
        CPUUsage:      getCPUUsage(),
        MemoryUsage:   m.Alloc,
        GoroutineCount: runtime.NumGoroutine(),
        AllocMemory:   m.Alloc,
        TotalMemory:   m.TotalAlloc,
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()
    
    time.Sleep(100 * time.Millisecond)
    
    end := time.Now()
    endGoroutines := runtime.NumGoroutine()
    
    duration := end.Sub(start).Seconds()
    goroutineChange := float64(endGoroutines - startGoroutines)
    
    return goroutineChange / duration * 100
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := collectMetrics()
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    if err := json.NewEncoder(w).Encode(metrics); err != nil {
        http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/metrics", metricsHandler)
    
    port := ":8080"
    fmt.Printf("Metrics collector server starting on port %s\n", port)
    fmt.Println("Access /metrics endpoint for system metrics")
    
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal("Server failed to start: ", err)
    }
}