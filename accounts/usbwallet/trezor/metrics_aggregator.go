package metrics

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	windowSize  time.Duration
	bucketSize  time.Duration
	buckets     []float64
	timestamps  []time.Time
	currentIdx  int
	mu          sync.RWMutex
}

func NewSlidingWindow(windowSize, bucketSize time.Duration) *SlidingWindow {
	numBuckets := int(windowSize / bucketSize)
	return &SlidingWindow{
		windowSize: windowSize,
		bucketSize: bucketSize,
		buckets:    make([]float64, numBuckets),
		timestamps: make([]time.Time, numBuckets),
	}
}

func (sw *SlidingWindow) Add(value float64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	sw.cleanupOldBuckets(now)

	sw.buckets[sw.currentIdx] += value
	sw.timestamps[sw.currentIdx] = now
}

func (sw *SlidingWindow) cleanupOldBuckets(now time.Time) {
	cutoff := now.Add(-sw.windowSize)

	for i := 0; i < len(sw.buckets); i++ {
		if sw.timestamps[i].Before(cutoff) {
			sw.buckets[i] = 0
			sw.timestamps[i] = time.Time{}
		}
	}

	if sw.timestamps[sw.currentIdx].IsZero() || 
	   now.Sub(sw.timestamps[sw.currentIdx]) >= sw.bucketSize {
		sw.currentIdx = (sw.currentIdx + 1) % len(sw.buckets)
		sw.buckets[sw.currentIdx] = 0
		sw.timestamps[sw.currentIdx] = time.Time{}
	}
}

func (sw *SlidingWindow) Sum() float64 {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	sw.cleanupOldBuckets(time.Now())
	
	var total float64
	for _, val := range sw.buckets {
		total += val
	}
	return total
}

func (sw *SlidingWindow) Average() float64 {
	sum := sw.Sum()
	
	sw.mu.RLock()
	validBuckets := 0
	for _, ts := range sw.timestamps {
		if !ts.IsZero() {
			validBuckets++
		}
	}
	sw.mu.RUnlock()

	if validBuckets == 0 {
		return 0
	}
	return sum / float64(validBuckets)
}