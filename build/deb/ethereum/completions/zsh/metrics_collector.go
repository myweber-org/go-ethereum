package main

import (
    "fmt"
    "log"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    GoroutineCount int
    MemoryAlloc    uint64
    MemoryTotal    uint64
    NumCPU         int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:      time.Now(),
        GoroutineCount: runtime.NumGoroutine(),
        MemoryAlloc:    m.Alloc,
        MemoryTotal:    m.TotalAlloc,
        NumCPU:         runtime.NumCPU(),
    }
}

func logMetrics(metrics SystemMetrics) {
    log.Printf(
        "Metrics collected at %s: Goroutines=%d, Alloc=%d bytes, TotalAlloc=%d bytes, CPUs=%d",
        metrics.Timestamp.Format(time.RFC3339),
        metrics.GoroutineCount,
        metrics.MemoryAlloc,
        metrics.MemoryTotal,
        metrics.NumCPU,
    )
}

func main() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    fmt.Println("Starting system metrics collector...")
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            logMetrics(metrics)
        }
    }
}