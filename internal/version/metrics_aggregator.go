package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize   time.Duration
	percentiles  []float64
	measurements []measurement
	mu           sync.RWMutex
}

type measurement struct {
	timestamp time.Time
	value     float64
}

func NewAggregator(windowSize time.Duration, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		percentiles: percentiles,
	}
}

func (a *Aggregator) Record(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.measurements = append(a.measurements, measurement{
		timestamp: now,
		value:     value,
	})
	a.cleanup(now)
}

func (a *Aggregator) GetStats() map[float64]float64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.cleanup(now)

	if len(a.measurements) == 0 {
		return make(map[float64]float64)
	}

	values := make([]float64, len(a.measurements))
	for i, m := range a.measurements {
		values[i] = m.value
	}
	sort.Float64s(values)

	stats := make(map[float64]float64)
	for _, p := range a.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		index := int(float64(len(values)-1) * p / 100.0)
		stats[p] = values[index]
	}

	return stats
}

func (a *Aggregator) cleanup(now time.Time) {
	cutoff := now.Add(-a.windowSize)
	validStart := 0
	for i, m := range a.measurements {
		if m.timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	a.measurements = a.measurements[validStart:]
}