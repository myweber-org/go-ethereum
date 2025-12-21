package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    Goroutines  int
    MemoryAlloc uint64
    CPUCores    int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now().UTC(),
        Goroutines:  runtime.NumGoroutine(),
        MemoryAlloc: m.Alloc,
        CPUCores:    runtime.NumCPU(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Printf("Memory Allocation: %d bytes\n", metrics.MemoryAlloc)
    fmt.Printf("CPU Cores: %d\n", metrics.CPUCores)
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