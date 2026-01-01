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
	windowSize  time.Duration
	metrics     []Metric
	mu          sync.RWMutex
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		metrics:     make([]Metric, 0),
		percentiles: percentiles,
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
	validStart := 0

	for i, metric := range swa.metrics {
		if metric.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}

	if validStart > 0 {
		swa.metrics = swa.metrics[validStart:]
	}
}

func (swa *SlidingWindowAggregator) CalculateStats() (map[string]float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	swa.cleanupOldMetrics(time.Now())

	if len(swa.metrics) == 0 {
		return nil, fmt.Errorf("no metrics in window")
	}

	values := make([]float64, len(swa.metrics))
	for i, metric := range swa.metrics {
		values[i] = metric.Value
	}

	sort.Float64s(values)

	stats := make(map[string]float64)
	stats["count"] = float64(len(values))
	stats["min"] = values[0]
	stats["max"] = values[len(values)-1]

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	stats["mean"] = sum / float64(len(values))

	if len(values)%2 == 0 {
		stats["median"] = (values[len(values)/2-1] + values[len(values)/2]) / 2
	} else {
		stats["median"] = values[len(values)/2]
	}

	variance := 0.0
	for _, v := range values {
		variance += math.Pow(v-stats["mean"], 2)
	}
	variance /= float64(len(values))
	stats["stddev"] = math.Sqrt(variance)

	for _, p := range swa.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		index := (p / 100) * float64(len(values)-1)
		lower := int(math.Floor(index))
		upper := int(math.Ceil(index))

		if lower == upper {
			stats[fmt.Sprintf("p%.0f", p)] = values[lower]
		} else {
			weight := index - float64(lower)
			stats[fmt.Sprintf("p%.0f", p)] = values[lower]*(1-weight) + values[upper]*weight
		}
	}

	return stats, nil
}

func (swa *SlidingWindowAggregator) GetCurrentMetrics() []Metric {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	swa.cleanupOldMetrics(time.Now())
	result := make([]Metric, len(swa.metrics))
	copy(result, swa.metrics)
	return result
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, []float64{50, 90, 95, 99})

	for i := 0; i < 100; i++ {
		value := float64(i) + math.Sin(float64(i)*0.1)*10
		aggregator.AddMetric(value)
		time.Sleep(100 * time.Millisecond)
	}

	stats, err := aggregator.CalculateStats()
	if err != nil {
		fmt.Printf("Error calculating stats: %v\n", err)
		return
	}

	fmt.Println("Aggregated Statistics:")
	for key, value := range stats {
		fmt.Printf("%s: %.4f\n", key, value)
	}

	fmt.Printf("\nTotal metrics in window: %d\n", len(aggregator.GetCurrentMetrics()))
}