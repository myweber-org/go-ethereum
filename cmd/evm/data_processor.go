package main

import (
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessRecords(records []Record) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Valid && r.Value > 0 {
			r.Name = strings.ToUpper(strings.TrimSpace(r.Name))
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func main() {
	records := []Record{
		{1, "alpha", 10.5, true},
		{2, "beta", -5.0, true},
		{3, "gamma", 7.2, false},
		{4, "delta", 3.8, true},
	}
	result := ProcessRecords(records)
	for _, r := range result {
		fmt.Printf("ID: %d, Name: %s, Value: %.1f\n", r.ID, r.Name, r.Value)
	}
}