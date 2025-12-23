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
    CPUPercent  float64
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now().UTC(),
        Goroutines:  runtime.NumGoroutine(),
        MemoryAlloc: m.Alloc,
        CPUPercent:  getCPUUsage(),
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start).Seconds()
    endGoroutines := runtime.NumGoroutine()

    return float64(endGoroutines-startGoroutines) / elapsed * 10
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            fmt.Printf("[%s] Goroutines: %d, Memory: %v bytes, CPU Load: %.2f%%\n",
                metrics.Timestamp.Format("15:04:05"),
                metrics.Goroutines,
                metrics.MemoryAlloc,
                metrics.CPUPercent)
        }
    }
}