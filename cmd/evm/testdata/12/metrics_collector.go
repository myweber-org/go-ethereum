package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    GoroutineCount int
    MemoryAlloc   uint64
    MemoryTotal   uint64
    NumCPU        int
}

func CollectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:     time.Now(),
        GoroutineCount: runtime.NumGoroutine(),
        MemoryAlloc:   memStats.Alloc,
        MemoryTotal:   memStats.TotalAlloc,
        NumCPU:        runtime.NumCPU(),
    }
}

func (m SystemMetrics) String() string {
    return fmt.Sprintf("[%s] Goroutines: %d, Memory: %s alloc / %s total, CPUs: %d",
        m.Timestamp.Format("15:04:05"),
        m.GoroutineCount,
        formatBytes(m.MemoryAlloc),
        formatBytes(m.MemoryTotal),
        m.NumCPU)
}

func formatBytes(bytes uint64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := uint64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := CollectMetrics()
            fmt.Println(metrics)
        }
    }
}