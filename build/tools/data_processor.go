
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
	rowCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csv read error: %w", err)
		}

		if !headerProcessed {
			headerProcessed = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write header: %w", err)
			}
			continue
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			if cleanedField == "" {
				cleanedField = "N/A"
			}
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
		rowCount++
	}

	fmt.Printf("Processed %d data rows successfully\n", rowCount)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV processing completed")
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func transformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	user.Username = sanitizeUsername(user.Username)

	if !validateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func processUserInput(input string) {
	rawData := []byte(input)
	user, err := transformUserData(rawData)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed user: ID=%d, Username='%s', Email='%s', Age=%d\n",
		user.ID, user.Username, user.Email, user.Age)
}

func main() {
	testData := `{"id": 1, "username": "  john_doe  ", "email": "john@example.com", "age": 30}`
	processUserInput(testData)

	invalidData := `{"id": 2, "username": "test", "email": "invalid-email", "age": 200}`
	processUserInput(invalidData)
}