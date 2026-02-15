package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Processed bool      `json:"processed"`
}

type Processor struct {
	mu          sync.RWMutex
	records     map[string]DataRecord
	workerCount int
}

func NewProcessor(workers int) *Processor {
	return &Processor{
		records:     make(map[string]DataRecord),
		workerCount: workers,
	}
}

func (p *Processor) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var records []DataRecord
	if err := json.Unmarshal(bytes, &records); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	p.mu.Lock()
	for _, record := range records {
		p.records[record.ID] = record
	}
	p.mu.Unlock()

	log.Printf("Loaded %d records from %s", len(records), filename)
	return nil
}

func (p *Processor) ProcessRecord(id string) error {
	p.mu.RLock()
	record, exists := p.records[id]
	p.mu.RUnlock()

	if !exists {
		return errors.New("record not found")
	}

	if record.Processed {
		return errors.New("record already processed")
	}

	time.Sleep(50 * time.Millisecond)

	record.Value = record.Value * 1.1
	record.Processed = true

	p.mu.Lock()
	p.records[id] = record
	p.mu.Unlock()

	return nil
}

func (p *Processor) RunConcurrentProcessing() {
	var wg sync.WaitGroup
	ids := p.getAllIDs()

	workChan := make(chan string, len(ids))

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for id := range workChan {
				if err := p.ProcessRecord(id); err != nil {
					log.Printf("Worker %d: Failed to process %s: %v", workerID, id, err)
				} else {
					log.Printf("Worker %d: Successfully processed %s", workerID, id)
				}
			}
		}(i)
	}

	for _, id := range ids {
		workChan <- id
	}
	close(workChan)

	wg.Wait()
}

func (p *Processor) getAllIDs() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ids := make([]string, 0, len(p.records))
	for id := range p.records {
		ids = append(ids, id)
	}
	return ids
}

func (p *Processor) SaveToFile(filename string) error {
	p.mu.RLock()
	records := make([]DataRecord, 0, len(p.records))
	for _, record := range p.records {
		records = append(records, record)
	}
	p.mu.RUnlock()

	bytes, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Saved %d records to %s", len(records), filename)
	return nil
}

func main() {
	processor := NewProcessor(4)

	if err := processor.LoadFromFile("input_data.json"); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	start := time.Now()
	processor.RunConcurrentProcessing()
	elapsed := time.Since(start)

	log.Printf("Processing completed in %v", elapsed)

	if err := processor.SaveToFile("output_data.json"); err != nil {
		log.Fatalf("Failed to save data: %v", err)
	}
}