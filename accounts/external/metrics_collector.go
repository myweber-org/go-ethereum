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
    GoroutineCount int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now(),
        MemoryUsedMB:  m.Alloc / 1024 / 1024,
        GoroutineCount: runtime.NumGoroutine(),
    }
}

func startMetricsCollector(interval time.Duration, stopChan <-chan struct{}) <-chan SystemMetrics {
    metricsChan := make(chan SystemMetrics)
    
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        defer close(metricsChan)
        
        for {
            select {
            case <-ticker.C:
                metricsChan <- collectMetrics()
            case <-stopChan:
                return
            }
        }
    }()
    
    return metricsChan
}

func main() {
    stopChan := make(chan struct{})
    defer close(stopChan)
    
    metricsStream := startMetricsCollector(5*time.Second, stopChan)
    
    for i := 0; i < 3; i++ {
        metrics := <-metricsStream
        fmt.Printf("Time: %v | Memory: %dMB | Goroutines: %d\n",
            metrics.Timestamp.Format("15:04:05"),
            metrics.MemoryUsedMB,
            metrics.GoroutineCount)
    }
}