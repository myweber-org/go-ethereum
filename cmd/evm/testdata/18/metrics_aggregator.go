
package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	mu sync.RWMutex
	latencies []time.Duration
	errors    int
	total     int
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		latencies: make([]time.Duration, 0),
	}
}

func (a *Aggregator) Record(latency time.Duration, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.latencies = append(a.latencies, latency)
	a.total++

	if err != nil {
		a.errors++
	}
}

func (a *Aggregator) GetStats() (avgLatency time.Duration, errorRate float64) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.total == 0 {
		return 0, 0
	}

	var sum time.Duration
	for _, l := range a.latencies {
		sum += l
	}
	avgLatency = sum / time.Duration(len(a.latencies))
	errorRate = float64(a.errors) / float64(a.total) * 100

	return avgLatency, errorRate
}

func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.latencies = make([]time.Duration, 0)
	a.errors = 0
	a.total = 0
}