package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Processed bool      `json:"processed"`
}

type Processor struct {
	mu          sync.RWMutex
	records     []DataRecord
	workerCount int
}

func NewProcessor(workerCount int) *Processor {
	return &Processor{
		records:     make([]DataRecord, 0),
		workerCount: workerCount,
	}
}

func (p *Processor) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var records []DataRecord
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	p.mu.Lock()
	p.records = records
	p.mu.Unlock()

	log.Printf("Loaded %d records from %s", len(records), filename)
	return nil
}

func (p *Processor) ProcessRecord(record DataRecord) DataRecord {
	time.Sleep(10 * time.Millisecond)
	record.Value = record.Value * 1.1
	record.Processed = true
	record.Timestamp = time.Now()
	return record
}

func (p *Processor) Run() []DataRecord {
	var wg sync.WaitGroup
	ch := make(chan DataRecord, len(p.records))
	results := make(chan DataRecord, len(p.records))

	p.mu.RLock()
	for _, record := range p.records {
		ch <- record
	}
	close(ch)
	p.mu.RUnlock()

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for record := range ch {
				log.Printf("Worker %d processing record ID %d", workerID, record.ID)
				processed := p.ProcessRecord(record)
				results <- processed
			}
		}(i)
	}

	wg.Wait()
	close(results)

	processedRecords := make([]DataRecord, 0)
	for result := range results {
		processedRecords = append(processedRecords, result)
	}

	return processedRecords
}

func (p *Processor) SaveToFile(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	log.Printf("Saved %d records to %s", len(records), filename)
	return nil
}

func generateSampleData() []DataRecord {
	records := make([]DataRecord, 100)
	for i := 0; i < 100; i++ {
		records[i] = DataRecord{
			ID:        i + 1,
			Name:      fmt.Sprintf("Record_%d", i+1),
			Value:     float64(i) * 1.5,
			Timestamp: time.Now(),
			Processed: false,
		}
	}
	return records
}

func main() {
	processor := NewProcessor(5)

	sampleData := generateSampleData()
	tempFile := "sample_data.json"

	file, err := os.Create(tempFile)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sampleData); err != nil {
		log.Fatal(err)
	}
	file.Close()
	defer os.Remove(tempFile)

	if err := processor.LoadFromFile(tempFile); err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	processed := processor.Run()
	elapsed := time.Since(startTime)

	log.Printf("Processed %d records in %v using %d workers", len(processed), elapsed, processor.workerCount)

	outputFile := "processed_data.json"
	if err := processor.SaveToFile(outputFile, processed); err != nil {
		log.Fatal(err)
	}

	log.Printf("Processing completed. Results saved to %s", outputFile)
}