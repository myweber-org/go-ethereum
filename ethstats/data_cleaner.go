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

func readCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
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

		if !isValidEmail(email) {
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

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func filterByScore(records []Record, minScore float64) []Record {
	filtered := []Record{}
	for _, record := range records {
		if record.Score >= minScore {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func writeCSV(filename string, records []Record) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Score"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			record.Email,
			fmt.Sprintf("%.2f", record.Score),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	records, err := readCSV("input.csv")
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Average score: %.2f\n", calculateAverage(records))

	filtered := filterByScore(records, 75.0)
	fmt.Printf("Records with score >= 75: %d\n", len(filtered))

	if err := writeCSV("output.csv", filtered); err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
		return
	}

	fmt.Println("Data processing completed successfully")
}