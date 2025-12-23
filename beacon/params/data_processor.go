package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

type Processor interface {
	Process(*DataRecord) error
}

type Validator struct{}

func (v Validator) Process(record *DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid ID")
	}
	if strings.TrimSpace(record.Name) == "" {
		return errors.New("empty name")
	}
	if record.Value < 0 {
		return errors.New("negative value")
	}
	record.Valid = true
	return nil
}

type Transformer struct{}

func (t Transformer) Process(record *DataRecord) error {
	if !record.Valid {
		return errors.New("cannot transform invalid record")
	}
	record.Name = strings.ToUpper(record.Name)
	record.Value = record.Value * 1.1
	return nil
}

type Pipeline struct {
	processors []Processor
}

func NewPipeline(processors ...Processor) *Pipeline {
	return &Pipeline{processors: processors}
}

func (p *Pipeline) Execute(record *DataRecord) error {
	for _, processor := range p.processors {
		if err := processor.Process(record); err != nil {
			return fmt.Errorf("pipeline error: %w", err)
		}
	}
	return nil
}

func main() {
	records := []DataRecord{
		{ID: 1, Name: "record_one", Value: 100.0},
		{ID: 2, Name: "  ", Value: 200.0},
		{ID: 0, Name: "record_three", Value: 300.0},
		{ID: 4, Name: "record_four", Value: -50.0},
	}

	pipeline := NewPipeline(Validator{}, Transformer{})

	for i := range records {
		err := pipeline.Execute(&records[i])
		if err != nil {
			fmt.Printf("Record %d failed: %v\n", records[i].ID, err)
		} else {
			fmt.Printf("Record %d processed: %s, %.2f\n", records[i].ID, records[i].Name, records[i].Value)
		}
	}
}