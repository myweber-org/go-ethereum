package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    CPUPercent    float64
    MemoryUsedMB  uint64
    MemoryTotalMB uint64
    GoroutineCount int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now(),
        MemoryUsedMB:  bToMb(m.Alloc),
        MemoryTotalMB: bToMb(m.Sys),
        GoroutineCount: runtime.NumGoroutine(),
        CPUPercent:    getCPUPercent(),
    }
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

func getCPUPercent() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()
    
    time.Sleep(100 * time.Millisecond)
    
    end := time.Now()
    endGoroutines := runtime.NumGoroutine()
    
    duration := end.Sub(start).Seconds()
    avgGoroutines := float64(startGoroutines+endGoroutines) / 2
    
    return avgGoroutines * duration * 0.1
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("Memory Used: %d MB\n", metrics.MemoryUsedMB)
    fmt.Printf("Total Memory: %d MB\n", metrics.MemoryTotalMB)
    fmt.Printf("Active Goroutines: %d\n", metrics.GoroutineCount)
    fmt.Printf("Memory Utilization: %.2f%%\n", 
        float64(metrics.MemoryUsedMB)/float64(metrics.MemoryTotalMB)*100)
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
            fmt.Println("---")
        }
    }
}