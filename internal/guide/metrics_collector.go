package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp    time.Time
	CPUPercent   float64
	MemoryAlloc  uint64
	MemoryTotal  uint64
	GoroutineNum int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:    time.Now(),
		CPUPercent:   getCPUUsage(),
		MemoryAlloc:  m.Alloc,
		MemoryTotal:  m.TotalAlloc,
		GoroutineNum: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	startCPU := runtime.NumCPU()
	time.Sleep(100 * time.Millisecond)
	end := time.Now()
	endCPU := runtime.NumCPU()

	elapsed := end.Sub(start).Seconds()
	cpuDelta := float64(endCPU - startCPU)

	if elapsed > 0 {
		return cpuDelta / elapsed * 100
	}
	return 0.0
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("Time: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
	fmt.Printf("Goroutines: %d\n", metrics.GoroutineNum)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics := collectMetrics()
			displayMetrics(metrics)
		}
	}
}