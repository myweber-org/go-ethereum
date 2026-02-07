package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    MemoryAlloc uint64
    MemoryTotal uint64
    NumCPU      int
    NumGoroutine int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:   time.Now(),
        MemoryAlloc: m.Alloc,
        MemoryTotal: m.TotalAlloc,
        NumCPU:      runtime.NumCPU(),
        NumGoroutine: runtime.NumGoroutine(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Metrics collected at: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
    fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
    fmt.Printf("CPU Cores: %d\n", metrics.NumCPU)
    fmt.Printf("Active Goroutines: %d\n", metrics.NumGoroutine)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}