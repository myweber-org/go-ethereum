package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Value int
	Valid bool
}

type Processor interface {
	Process(DataRecord) (DataRecord, error)
}

type Validator struct{}

func (v Validator) Process(rec DataRecord) (DataRecord, error) {
	if rec.ID == "" {
		return rec, errors.New("empty ID")
	}
	if rec.Value <= 0 {
		return rec, errors.New("invalid value")
	}
	rec.Valid = true
	return rec, nil
}

type Transformer struct{}

func (t Transformer) Process(rec DataRecord) (DataRecord, error) {
	if !rec.Valid {
		return rec, errors.New("cannot transform invalid record")
	}
	rec.ID = strings.ToUpper(rec.ID)
	rec.Value = rec.Value * 2
	return rec, nil
}

type Pipeline struct {
	processors []Processor
}

func (p *Pipeline) AddProcessor(proc Processor) {
	p.processors = append(p.processors, proc)
}

func (p *Pipeline) Execute(rec DataRecord) (DataRecord, error) {
	var err error
	for _, proc := range p.processors {
		rec, err = proc.Process(rec)
		if err != nil {
			return rec, fmt.Errorf("pipeline error: %w", err)
		}
	}
	return rec, nil
}

func main() {
	pipeline := &Pipeline{}
	pipeline.AddProcessor(Validator{})
	pipeline.AddProcessor(Transformer{})

	records := []DataRecord{
		{"abc123", 42, false},
		{"", 100, false},
		{"def456", -5, false},
	}

	for _, rec := range records {
		result, err := pipeline.Execute(rec)
		if err != nil {
			fmt.Printf("Error processing %v: %v\n", rec.ID, err)
			continue
		}
		fmt.Printf("Processed: ID=%s, Value=%d, Valid=%v\n",
			result.ID, result.Value, result.Valid)
	}
}