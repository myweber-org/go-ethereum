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
}