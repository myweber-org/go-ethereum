
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataCleaner struct {
    inputPath  string
    outputPath string
    delimiter  rune
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        inputPath:  input,
        outputPath: output,
        delimiter:  ',',
    }
}

func (dc *DataCleaner) SetDelimiter(delim rune) {
    dc.delimiter = delim
}

func (dc *DataCleaner) RemoveDuplicates() error {
    inFile, err := os.Open(dc.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(dc.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    reader.Comma = dc.delimiter
    writer := csv.NewWriter(outFile)
    writer.Comma = dc.delimiter
    defer writer.Flush()

    seen := make(map[string]bool)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read CSV record: %w", err)
        }

        key := strings.Join(record, "|")
        if !seen[key] {
            seen[key] = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("failed to write CSV record: %w", err)
            }
        }
    }

    return nil
}

func main() {
    cleaner := NewDataCleaner("input.csv", "output.csv")
    cleaner.SetDelimiter(';')
    
    if err := cleaner.RemoveDuplicates(); err != nil {
        fmt.Printf("Error cleaning data: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Data cleaning completed successfully")
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	cleaned := dc.CleanString(value)
	if dc.seen[cleaned] {
		return true
	}
	dc.seen[cleaned] = true
	return false
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	cleaned := dc.CleanString(email)
	return strings.Contains(cleaned, "@") && strings.Contains(cleaned, ".")
}

func (dc *DataCleaner) GetUniqueCount() int {
	return len(dc.seen)
}

func main() {
	cleaner := NewDataCleaner()

	emails := []string{
		"  USER@EXAMPLE.COM  ",
		"user@example.com",
		"test@domain.org",
		"invalid-email",
		"TEST@DOMAIN.ORG  ",
	}

	fmt.Println("Processing emails:")
	for _, email := range emails {
		cleaned := cleaner.CleanString(email)
		isDup := cleaner.IsDuplicate(email)
		valid := cleaner.ValidateEmail(email)

		fmt.Printf("Original: %q\n", email)
		fmt.Printf("Cleaned: %q\n", cleaned)
		fmt.Printf("Duplicate: %v\n", isDup)
		fmt.Printf("Valid: %v\n\n", valid)
	}

	fmt.Printf("Total unique entries: %d\n", cleaner.GetUniqueCount())
}