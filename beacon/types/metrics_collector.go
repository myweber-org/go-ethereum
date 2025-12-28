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

    metrics := SystemMetrics{
        Timestamp:     time.Now(),
        MemoryUsedMB:  m.Alloc / 1024 / 1024,
        MemoryTotalMB: m.Sys / 1024 / 1024,
        GoroutineCount: runtime.NumGoroutine(),
    }

    metrics.CPUPercent = calculateCPUUsage()
    return metrics
}

func calculateCPUUsage() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start).Seconds()
    endGoroutines := runtime.NumGoroutine()

    usage := (float64(endGoroutines-startGoroutines) * 0.1) + (elapsed * 5.0)
    if usage > 100.0 {
        usage = 100.0
    }
    if usage < 0.0 {
        usage = 0.0
    }
    return usage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Metrics at %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("  CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("  Memory: %dMB / %dMB\n", metrics.MemoryUsedMB, metrics.MemoryTotalMB)
    fmt.Printf("  Goroutines: %d\n", metrics.GoroutineCount)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for i := 0; i < 3; i++ {
        metrics := collectMetrics()
        printMetrics(metrics)
        <-ticker.C
    }
}