package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemoryAlloc uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:   time.Now().UTC(),
        CPUPercent:  getCPUUsage(),
        MemoryAlloc: m.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    var startStats runtime.MemStats
    runtime.ReadMemStats(&startStats)
    
    time.Sleep(100 * time.Millisecond)
    
    end := time.Now()
    var endStats runtime.MemStats
    runtime.ReadMemStats(&endStats)
    
    elapsed := end.Sub(start).Seconds()
    if elapsed == 0 {
        return 0.0
    }
    
    cpuUsage := float64(endStats.NumGC-startStats.NumGC) * 100 / elapsed
    if cpuUsage < 0 {
        cpuUsage = 0
    }
    return cpuUsage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] CPU: %.2f%% | Memory: %v MB | Goroutines: %d\n",
        metrics.Timestamp.Format("2006-01-02 15:04:05"),
        metrics.CPUPercent,
        metrics.MemoryAlloc/1024/1024,
        metrics.Goroutines)
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