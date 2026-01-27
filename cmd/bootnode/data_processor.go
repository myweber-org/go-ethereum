
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

type Validator interface {
	Validate(DataRecord) error
}

type Transformer interface {
	Transform(DataRecord) DataRecord
}

type LengthValidator struct {
	MinLength int
}

func (v LengthValidator) Validate(record DataRecord) error {
	if len(record.Name) < v.MinLength {
		return errors.New("name too short")
	}
	return nil
}

type ValueValidator struct {
	MinValue float64
	MaxValue float64
}

func (v ValueValidator) Validate(record DataRecord) error {
	if record.Value < v.MinValue || record.Value > v.MaxValue {
		return errors.New("value out of range")
	}
	return nil
}

type UppercaseTransformer struct{}

func (t UppercaseTransformer) Transform(record DataRecord) DataRecord {
	record.Name = strings.ToUpper(record.Name)
	return record
}

type MultiplierTransformer struct {
	Factor float64
}

func (t MultiplierTransformer) Transform(record DataRecord) DataRecord {
	record.Value = record.Value * t.Factor
	return record
}

type ProcessingPipeline struct {
	validators   []Validator
	transformers []Transformer
}

func NewProcessingPipeline() *ProcessingPipeline {
	return &ProcessingPipeline{
		validators:   make([]Validator, 0),
		transformers: make([]Transformer, 0),
	}
}

func (p *ProcessingPipeline) AddValidator(v Validator) {
	p.validators = append(p.validators, v)
}

func (p *ProcessingPipeline) AddTransformer(t Transformer) {
	p.transformers = append(p.transformers, t)
}

func (p *ProcessingPipeline) Process(record DataRecord) (DataRecord, error) {
	for _, validator := range p.validators {
		if err := validator.Validate(record); err != nil {
			record.Valid = false
			return record, fmt.Errorf("validation failed: %w", err)
		}
	}

	processedRecord := record
	for _, transformer := range p.transformers {
		processedRecord = transformer.Transform(processedRecord)
	}

	processedRecord.Valid = true
	return processedRecord, nil
}

func main() {
	pipeline := NewProcessingPipeline()
	pipeline.AddValidator(LengthValidator{MinLength: 3})
	pipeline.AddValidator(ValueValidator{MinValue: 0, MaxValue: 100})
	pipeline.AddTransformer(UppercaseTransformer{})
	pipeline.AddTransformer(MultiplierTransformer{Factor: 1.5})

	records := []DataRecord{
		{ID: 1, Name: "test", Value: 50.0},
		{ID: 2, Name: "ab", Value: 75.0},
		{ID: 3, Name: "sample", Value: 150.0},
	}

	for _, record := range records {
		result, err := pipeline.Process(record)
		if err != nil {
			fmt.Printf("Record %d processing failed: %v\n", record.ID, err)
		} else {
			fmt.Printf("Processed Record %d: %+v\n", result.ID, result)
		}
	}
}