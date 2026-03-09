package main

import (
    "errors"
    "strings"
    "unicode"
)

type UserProfile struct {
    Username string
    Email    string
    Age      int
}

func ValidateProfile(p UserProfile) error {
    if strings.TrimSpace(p.Username) == "" {
        return errors.New("username cannot be empty")
    }
    if len(p.Username) < 3 || len(p.Username) > 20 {
        return errors.New("username must be between 3 and 20 characters")
    }
    for _, r := range p.Username {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
            return errors.New("username can only contain letters, digits, and underscores")
        }
    }

    if !strings.Contains(p.Email, "@") {
        return errors.New("invalid email format")
    }

    if p.Age < 0 || p.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }

    return nil
}

func NormalizeProfile(p UserProfile) UserProfile {
    normalized := p
    normalized.Username = strings.ToLower(strings.TrimSpace(p.Username))
    normalized.Email = strings.ToLower(strings.TrimSpace(p.Email))
    return normalized
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
    if err := ValidateProfile(p); err != nil {
        return UserProfile{}, err
    }
    return NormalizeProfile(p), nil
}package main

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	pattern := `^[a-zA-Z0-9\s\.\-_]+$`
	matched, err := regexp.MatchString(pattern, trimmed)
	if err != nil || !matched {
		return "", false
	}

	return trimmed, true
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []DataRecord
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        records = append(records, record)
    }

    return records, nil
}

func ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
        return false
    }
    return strings.Contains(parts[1], ".")
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
    var activeUsers []DataRecord
    for _, record := range records {
        if strings.ToLower(record.Active) == "true" || record.Active == "1" {
            if ValidateEmail(record.Email) {
                activeUsers = append(activeUsers, record)
            }
        }
    }
    return activeUsers
}

func GenerateReport(records []DataRecord) {
    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Println("Active users with valid emails:")
    for i, record := range records {
        fmt.Printf("%d. ID: %s, Name: %s, Email: %s\n", i+1, record.ID, record.Name, record.Email)
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_processor.go <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := ProcessCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    activeUsers := FilterActiveUsers(records)
    GenerateReport(activeUsers)
}