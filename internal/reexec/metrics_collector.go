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
	NumGoroutine int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:   time.Now(),
		MemoryAlloc: m.Alloc,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("[%s] Memory: %v bytes, Goroutines: %d\n",
		metrics.Timestamp.Format("15:04:05"),
		metrics.MemoryAlloc,
		metrics.NumGoroutine)
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := collectMetrics()
		displayMetrics(metrics)
	}
}