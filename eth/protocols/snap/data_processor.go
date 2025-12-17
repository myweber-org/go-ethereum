
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type Record struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func readCSV(filename string, out chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading CSV: %v\n", err)
			continue
		}
		out <- row
	}
}

func processRow(row []string, out chan<- Record, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(row) != 4 {
		return
	}

	id, err := strconv.Atoi(row[0])
	if err != nil {
		return
	}

	name := row[1]

	value, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return
	}

	valid, err := strconv.ParseBool(row[3])
	if err != nil {
		return
	}

	out <- Record{
		ID:    id,
		Name:  name,
		Value: value,
		Valid: valid,
	}
}

func filterRecords(in <-chan Record, out chan<- Record, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	for record := range in {
		if record.Valid && record.Value > 100.0 {
			out <- record
		}
	}
}

func main() {
	csvRows := make(chan []string, 10)
	rawRecords := make(chan Record, 10)
	filteredRecords := make(chan Record, 10)

	var wg sync.WaitGroup

	wg.Add(1)
	go readCSV("input.csv", csvRows, &wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(rawRecords)

		var processWg sync.WaitGroup
		for row := range csvRows {
			processWg.Add(1)
			go processRow(row, rawRecords, &processWg)
		}
		processWg.Wait()
	}()

	wg.Add(1)
	go filterRecords(rawRecords, filteredRecords, &wg)

	go func() {
		wg.Wait()
		close(filteredRecords)
	}()

	var results []Record
	for record := range filteredRecords {
		results = append(results, record)
		fmt.Printf("Processed: ID=%d, Name=%s, Value=%.2f\n", record.ID, record.Name, record.Value)
	}

	fmt.Printf("Total filtered records: %d\n", len(results))
}