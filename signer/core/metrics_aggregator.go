package metrics

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	windowSize  time.Duration
	maxSamples  int
	samples     []float64
	timestamps  []time.Time
	mu          sync.RWMutex
	aggregation map[string]float64
}

func NewSlidingWindow(windowSize time.Duration, maxSamples int) *SlidingWindow {
	return &SlidingWindow{
		windowSize:  windowSize,
		maxSamples:  maxSamples,
		samples:     make([]float64, 0, maxSamples),
		timestamps:  make([]time.Time, 0, maxSamples),
		aggregation: make(map[string]float64),
	}
}

func (sw *SlidingWindow) AddSample(value float64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	sw.samples = append(sw.samples, value)
	sw.timestamps = append(sw.timestamps, now)

	sw.pruneOldSamples(now)
	sw.updateAggregation()
}

func (sw *SlidingWindow) pruneOldSamples(currentTime time.Time) {
	cutoff := currentTime.Add(-sw.windowSize)
	validStart := 0

	for i, ts := range sw.timestamps {
		if ts.After(cutoff) {
			validStart = i
			break
		}
	}

	if validStart > 0 {
		sw.samples = sw.samples[validStart:]
		sw.timestamps = sw.timestamps[validStart:]
	}

	if len(sw.samples) > sw.maxSamples {
		overflow := len(sw.samples) - sw.maxSamples
		sw.samples = sw.samples[overflow:]
		sw.timestamps = sw.timestamps[overflow:]
	}
}

func (sw *SlidingWindow) updateAggregation() {
	if len(sw.samples) == 0 {
		sw.aggregation = make(map[string]float64)
		return
	}

	var sum, min, max float64
	min = sw.samples[0]
	max = sw.samples[0]

	for _, v := range sw.samples {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	sw.aggregation["sum"] = sum
	sw.aggregation["avg"] = sum / float64(len(sw.samples))
	sw.aggregation["min"] = min
	sw.aggregation["max"] = max
	sw.aggregation["count"] = float64(len(sw.samples))
}

func (sw *SlidingWindow) GetAggregation() map[string]float64 {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	result := make(map[string]float64)
	for k, v := range sw.aggregation {
		result[k] = v
	}
	return result
}

func (sw *SlidingWindow) GetCurrentWindowSize() int {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return len(sw.samples)
}