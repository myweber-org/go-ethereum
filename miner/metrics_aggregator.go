package metrics

import (
	"sync"
	"time"
)

type LatencyAggregator struct {
	mu          sync.RWMutex
	count       int64
	total       time.Duration
	min         time.Duration
	max         time.Duration
	percentiles map[int]time.Duration
}

func NewLatencyAggregator() *LatencyAggregator {
	return &LatencyAggregator{
		min:         time.Hour,
		percentiles: make(map[int]time.Duration),
	}
}

func (la *LatencyAggregator) Record(duration time.Duration) {
	la.mu.Lock()
	defer la.mu.Unlock()

	la.count++
	la.total += duration

	if duration < la.min {
		la.min = duration
	}
	if duration > la.max {
		la.max = duration
	}
}

func (la *LatencyAggregator) Average() time.Duration {
	la.mu.RLock()
	defer la.mu.RUnlock()

	if la.count == 0 {
		return 0
	}
	return la.total / time.Duration(la.count)
}

func (la *LatencyAggregator) Stats() (count int64, avg, min, max time.Duration) {
	la.mu.RLock()
	defer la.mu.RUnlock()

	if la.count == 0 {
		return 0, 0, 0, 0
	}
	return la.count, la.total / time.Duration(la.count), la.min, la.max
}

func (la *LatencyAggregator) Reset() {
	la.mu.Lock()
	defer la.mu.Unlock()

	la.count = 0
	la.total = 0
	la.min = time.Hour
	la.max = 0
	la.percentiles = make(map[int]time.Duration)
}