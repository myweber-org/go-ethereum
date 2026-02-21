package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID    int
	Name  string
	Email string
	Score float64
}

func parseCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lineNum++

		if lineNum == 1 {
			continue
		}

		if len(line) != 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(line[1])
		email := strings.TrimSpace(line[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
		if err != nil {
			continue
		}

		if !validateEmail(email) {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		})
	}

	return records, nil
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	for _, r := range records {
		total += r.Score
	}
	return total / float64(len(records))
}

func filterByScore(records []Record, threshold float64) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Score >= threshold {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func main() {
	records, err := parseCSV("data.csv")
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	fmt.Printf("Total valid records: %d\n", len(records))
	fmt.Printf("Average score: %.2f\n", calculateAverage(records))

	highScorers := filterByScore(records, 80.0)
	fmt.Printf("Records with score >= 80: %d\n", len(highScorers))

	for _, r := range highScorers {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Score: %.1f\n",
			r.ID, r.Name, r.Email, r.Score)
	}
}