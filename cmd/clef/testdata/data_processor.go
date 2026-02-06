
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func processCSVFile(inputPath, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
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

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = strings.TrimSpace(field)
            if cleanedRecord[i] == "" {
                cleanedRecord[i] = "N/A"
            }
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run data_processor.go <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %s to %s\n", inputFile, outputFile)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25,"active":true}`)
	
	processedData, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed User Data:\n")
	fmt.Printf("Email: %s\n", processedData.Email)
	fmt.Printf("Username: %s\n", processedData.Username)
	fmt.Printf("Age: %d\n", processedData.Age)
	fmt.Printf("Active: %t\n", processedData.Active)
}