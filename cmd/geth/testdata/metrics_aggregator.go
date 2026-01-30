
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize  time.Duration
	maxSamples  int
	mu          sync.RWMutex
	samples     []float64
	timestamps  []time.Time
	percentiles []float64
}

func NewAggregator(windowSize time.Duration, maxSamples int, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		maxSamples:  maxSamples,
		percentiles: percentiles,
		samples:     make([]float64, 0, maxSamples),
		timestamps:  make([]time.Time, 0, maxSamples),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.samples = append(a.samples, value)
	a.timestamps = append(a.timestamps, now)

	a.evictOldSamples(now)
	if len(a.samples) > a.maxSamples {
		a.samples = a.samples[1:]
		a.timestamps = a.timestamps[1:]
	}
}

func (a *Aggregator) evictOldSamples(cutoff time.Time) {
	threshold := cutoff.Add(-a.windowSize)
	firstValid := 0
	for i, ts := range a.timestamps {
		if ts.After(threshold) {
			firstValid = i
			break
		}
	}
	if firstValid > 0 {
		a.samples = a.samples[firstValid:]
		a.timestamps = a.timestamps[firstValid:]
	}
}

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.samples) == 0 {
		return nil
	}

	stats := make(map[string]float64)
	sorted := make([]float64, len(a.samples))
	copy(sorted, a.samples)
	sort.Float64s(sorted)

	stats["count"] = float64(len(sorted))
	stats["min"] = sorted[0]
	stats["max"] = sorted[len(sorted)-1]

	var sum float64
	for _, v := range sorted {
		sum += v
	}
	stats["mean"] = sum / float64(len(sorted))

	for _, p := range a.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		idx := int(float64(len(sorted)-1) * p / 100.0)
		key := "p" + formatPercentile(p)
		stats[key] = sorted[idx]
	}

	return stats
}

func formatPercentile(p float64) string {
	if p == float64(int(p)) {
		return string(rune('0' + int(p)/10)) + string(rune('0' + int(p)%10))
	}
	return "0"
}