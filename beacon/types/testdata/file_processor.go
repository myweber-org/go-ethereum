
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int
	Content   string
	Valid     bool
	Timestamp time.Time
}

type Processor struct {
	records []DataRecord
	mu      sync.RWMutex
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(content string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        len(p.records) + 1,
		Content:   content,
		Valid:     len(content) > 0,
		Timestamp: time.Now(),
	}

	p.records = append(p.records, record)
}

func (p *Processor) ValidateRecords() (int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.records) == 0 {
		return 0, errors.New("no records to validate")
	}

	validCount := 0
	for _, record := range p.records {
		if record.Valid {
			validCount++
		}
	}

	return validCount, nil
}

func (p *Processor) ProcessBatch(contents []string) {
	var wg sync.WaitGroup
	results := make(chan DataRecord, len(contents))

	for _, content := range contents {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			p.AddRecord(c)
			p.mu.RLock()
			lastRecord := p.records[len(p.records)-1]
			p.mu.RUnlock()
			results <- lastRecord
		}(content)
	}

	wg.Wait()
	close(results)

	fmt.Println("Processing completed:")
	for record := range results {
		fmt.Printf("Record %d: %s (Valid: %v)\n",
			record.ID, record.Content, record.Valid)
	}
}

func main() {
	processor := NewProcessor()

	sampleData := []string{
		"alpha",
		"",
		"beta",
		"gamma",
		"",
		"delta",
	}

	processor.ProcessBatch(sampleData)

	validCount, err := processor.ValidateRecords()
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	fmt.Printf("Total valid records: %d\n", validCount)
}