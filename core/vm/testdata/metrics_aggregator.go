package aggregator

import (
	"sync"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxSamples  int
	metrics     []Metric
	mu          sync.RWMutex
	subscribers []chan float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxSamples int) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize: windowSize,
		maxSamples: maxSamples,
		metrics:    make([]Metric, 0, maxSamples),
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	metric := Metric{
		Timestamp: now,
		Value:     value,
	}

	swa.metrics = append(swa.metrics, metric)
	swa.cleanupOldMetrics(now)

	if len(swa.subscribers) > 0 {
		avg := swa.calculateAverage()
		for _, ch := range swa.subscribers {
			select {
			case ch <- avg:
			default:
			}
		}
	}
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

	swa.metrics = swa.metrics[validStart:]

	if len(swa.metrics) > swa.maxSamples {
		swa.metrics = swa.metrics[len(swa.metrics)-swa.maxSamples:]
	}
}

func (swa *SlidingWindowAggregator) calculateAverage() float64 {
	if len(swa.metrics) == 0 {
		return 0
	}

	var sum float64
	for _, metric := range swa.metrics {
		sum += metric.Value
	}
	return sum / float64(len(swa.metrics))
}

func (swa *SlidingWindowAggregator) GetAverage() float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()
	return swa.calculateAverage()
}

func (swa *SlidingWindowAggregator) Subscribe() <-chan float64 {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	ch := make(chan float64, 1)
	swa.subscribers = append(swa.subscribers, ch)
	return ch
}

func (swa *SlidingWindowAggregator) Close() {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	for _, ch := range swa.subscribers {
		close(ch)
	}
	swa.subscribers = nil
}