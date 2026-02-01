
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func removeDuplicates(inputFile, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	seen := make(map[string]bool)
	reader := csv.NewReader(in)
	var records [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%v", record)
		if !seen[key] {
			seen[key] = true
			records = append(records, record)
		}
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	return writer.WriteAll(records)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := removeDuplicates(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Duplicate removal completed successfully")
}