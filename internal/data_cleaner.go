
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func removeDuplicates(records [][]string) [][]string {
	seen := make(map[string]bool)
	var unique [][]string

	for _, record := range records {
		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func filterValidEmails(records [][]string, emailColumn int) [][]string {
	var valid [][]string
	valid = append(valid, records[0])

	for i := 1; i < len(records); i++ {
		if len(records[i]) > emailColumn && validateEmail(records[i][emailColumn]) {
			valid = append(valid, records[i])
		}
	}
	return valid
}

func sortByColumn(records [][]string, column int) {
	if len(records) < 2 {
		return
	}
	header := records[0]
	data := records[1:]

	sort.Slice(data, func(i, j int) bool {
		if len(data[i]) <= column || len(data[j]) <= column {
			return false
		}
		return data[i][column] < data[j][column]
	})

	records = append([][]string{header}, data...)
}

func writeCSV(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(records)
}

func processData(inputFile, outputFile string) error {
	records, err := readCSV(inputFile)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("empty file")
	}

	records = removeDuplicates(records)
	records = filterValidEmails(records, 2)
	sortByColumn(records, 1)

	return writeCSV(outputFile, records)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := processData(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data processing completed successfully")
}