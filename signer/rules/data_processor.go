
package data

import (
	"errors"
	"strings"
	"time"
)

type Record struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if r.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if r.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecords(records []Record, multiplier float64) []Record {
	transformed := make([]Record, len(records))
	for i, r := range records {
		transformed[i] = Record{
			ID:        strings.ToUpper(r.ID),
			Value:     r.Value * multiplier,
			Timestamp: r.Timestamp.UTC(),
			Tags:      append([]string{}, r.Tags...),
		}
	}
	return transformed
}

func FilterByThreshold(records []Record, threshold float64) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Value >= threshold {
			filtered = append(filtered, r)
		}
	}
	return filtered
}