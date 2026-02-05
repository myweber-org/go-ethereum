
package metrics

import (
	"sync"
	"time"
)

type RequestMetrics struct {
	mu               sync.RWMutex
	TotalRequests    int64
	Status2xx        int64
	Status4xx        int64
	Status5xx        int64
	TotalLatency     time.Duration
	MaxLatency       time.Duration
	MinLatency       time.Duration
	latencySamples   []time.Duration
	sampleCapacity   int
}

func NewRequestMetrics(sampleCapacity int) *RequestMetrics {
	return &RequestMetrics{
		latencySamples: make([]time.Duration, 0, sampleCapacity),
		sampleCapacity: sampleCapacity,
		MinLatency:     time.Hour,
	}
}

func (rm *RequestMetrics) RecordRequest(status int, latency time.Duration) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.TotalRequests++
	rm.TotalLatency += latency

	if latency > rm.MaxLatency {
		rm.MaxLatency = latency
	}
	if latency < rm.MinLatency {
		rm.MinLatency = latency
	}

	if len(rm.latencySamples) < rm.sampleCapacity {
		rm.latencySamples = append(rm.latencySamples, latency)
	} else {
		rm.latencySamples = append(rm.latencySamples[1:], latency)
	}

	switch {
	case status >= 200 && status < 300:
		rm.Status2xx++
	case status >= 400 && status < 500:
		rm.Status4xx++
	case status >= 500:
		rm.Status5xx++
	}
}

func (rm *RequestMetrics) AverageLatency() time.Duration {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if rm.TotalRequests == 0 {
		return 0
	}
	return rm.TotalLatency / time.Duration(rm.TotalRequests)
}

func (rm *RequestMetrics) LatencyPercentile(p float64) time.Duration {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if len(rm.latencySamples) == 0 {
		return 0
	}

	samples := make([]time.Duration, len(rm.latencySamples))
	copy(samples, rm.latencySamples)

	for i := 0; i < len(samples); i++ {
		for j := i + 1; j < len(samples); j++ {
			if samples[i] > samples[j] {
				samples[i], samples[j] = samples[j], samples[i]
			}
		}
	}

	index := int(float64(len(samples)) * p / 100.0)
	if index >= len(samples) {
		index = len(samples) - 1
	}
	return samples[index]
}

func (rm *RequestMetrics) Reset() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.TotalRequests = 0
	rm.Status2xx = 0
	rm.Status4xx = 0
	rm.Status5xx = 0
	rm.TotalLatency = 0
	rm.MaxLatency = 0
	rm.MinLatency = time.Hour
	rm.latencySamples = rm.latencySamples[:0]
}