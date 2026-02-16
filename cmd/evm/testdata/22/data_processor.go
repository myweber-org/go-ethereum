
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headerProcessed := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		if !headerProcessed {
			headerProcessed = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
			continue
		}

		cleanedRecord := dp.cleanRecord(record)
		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("error writing record: %w", err)
			}
		}
	}

	return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data processing completed successfully")
}package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Config struct {
    ServerAddress string `json:"server_address"`
    Port          int    `json:"port"`
    EnableLogging bool   `json:"enable_logging"`
    MaxConnections int   `json:"max_connections"`
}

func LoadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("failed to decode JSON: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.ServerAddress == "" {
        return fmt.Errorf("server_address cannot be empty")
    }
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    if c.MaxConnections < 1 {
        return fmt.Errorf("max_connections must be at least 1")
    }
    return nil
}