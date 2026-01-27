
package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	windowSize time.Duration
	metrics    []Metric
	mu         sync.RWMutex
}

func NewSlidingWindowAggregator(windowSize time.Duration) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize: windowSize,
		metrics:    make([]Metric, 0),
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics = append(swa.metrics, Metric{Timestamp: now, Value: value})
	swa.cleanupOldMetrics(now)
}

func (swa *SlidingWindowAggregator) cleanupOldMetrics(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	i := 0
	for i < len(swa.metrics) && swa.metrics[i].Timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		swa.metrics = swa.metrics[i:]
	}
}

func (swa *SlidingWindowAggregator) CalculatePercentile(p float64) (float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	swa.cleanupOldMetrics(time.Now())

	if len(swa.metrics) == 0 {
		return 0, fmt.Errorf("no metrics available")
	}

	values := make([]float64, len(swa.metrics))
	for i, m := range swa.metrics {
		values[i] = m.Value
	}
	sort.Float64s(values)

	index := p * float64(len(values)-1) / 100.0
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return values[lower], nil
	}

	weight := index - float64(lower)
	return values[lower]*(1-weight) + values[upper]*weight, nil
}

func (swa *SlidingWindowAggregator) GetStats() (float64, float64, float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	swa.cleanupOldMetrics(time.Now())

	if len(swa.metrics) == 0 {
		return 0, 0, 0, fmt.Errorf("no metrics available")
	}

	var sum float64
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, m := range swa.metrics {
		sum += m.Value
		if m.Value < min {
			min = m.Value
		}
		if m.Value > max {
			max = m.Value
		}
	}

	avg := sum / float64(len(swa.metrics))
	return avg, min, max, nil
}

func main() {
	aggregator := NewSlidingWindowAggregator(5 * time.Minute)

	for i := 0; i < 100; i++ {
		value := float64(i) + (float64(i%10) * 0.1)
		aggregator.AddMetric(value)
		time.Sleep(100 * time.Millisecond)
	}

	avg, min, max, err := aggregator.GetStats()
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
		return
	}
	fmt.Printf("Average: %.2f, Min: %.2f, Max: %.2f\n", avg, min, max)

	p95, err := aggregator.CalculatePercentile(95)
	if err != nil {
		fmt.Printf("Error calculating percentile: %v\n", err)
		return
	}
	fmt.Printf("95th Percentile: %.2f\n", p95)
}