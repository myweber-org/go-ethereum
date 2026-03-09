
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int
	Value     string
	Timestamp time.Time
	Valid     bool
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

func (p *Processor) AddRecord(id int, value string) error {
	if id <= 0 {
		return errors.New("invalid record ID")
	}
	if value == "" {
		return errors.New("record value cannot be empty")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        id,
		Value:     value,
		Timestamp: time.Now().UTC(),
		Valid:     true,
	}

	p.records = append(p.records, record)
	return nil
}

func (p *Processor) ValidateRecords() {
	p.mu.RLock()
	records := make([]DataRecord, len(p.records))
	copy(records, p.records)
	p.mu.RUnlock()

	var wg sync.WaitGroup
	for i := range records {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p.validateRecord(&records[idx])
		}(i)
	}
	wg.Wait()

	p.mu.Lock()
	p.records = records
	p.mu.Unlock()
}

func (p *Processor) validateRecord(record *DataRecord) {
	if record.ID <= 0 || record.Value == "" {
		record.Valid = false
		return
	}

	if time.Since(record.Timestamp) > 24*time.Hour {
		record.Valid = false
		return
	}

	record.Valid = true
}

func (p *Processor) GetValidCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, record := range p.records {
		if record.Valid {
			count++
		}
	}
	return count
}

func (p *Processor) DisplaySummary() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	fmt.Printf("Total records: %d\n", len(p.records))
	fmt.Printf("Valid records: %d\n", p.GetValidCount())
	fmt.Println("Record details:")
	for _, record := range p.records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("  ID: %d, Value: %s, Time: %s, Status: %s\n",
			record.ID, record.Value, record.Timestamp.Format(time.RFC3339), status)
	}
}

func main() {
	processor := NewProcessor()

	records := []struct {
		id    int
		value string
	}{
		{1, "alpha"},
		{2, "beta"},
		{3, ""},
		{0, "gamma"},
		{4, "delta"},
	}

	for _, r := range records {
		if err := processor.AddRecord(r.id, r.value); err != nil {
			fmt.Printf("Failed to add record %d: %v\n", r.id, err)
		}
	}

	processor.ValidateRecords()
	processor.DisplaySummary()
}