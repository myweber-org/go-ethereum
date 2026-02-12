package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func parseLogLine(line string) (*LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return nil, fmt.Errorf("invalid log format")
	}

	return &LogEntry{
		Timestamp: matches[1],
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[3],
	}, nil
}

func filterErrors(entries []LogEntry) []LogEntry {
	var errors []LogEntry
	for _, entry := range entries {
		if entry.Level == "ERROR" {
			errors = append(errors, entry)
		}
	}
	return errors
}

func processLogFile(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err == nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	entries, err := processLogFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	errorEntries := filterErrors(entries)
	fmt.Printf("Found %d error entries:\n", len(errorEntries))
	for _, entry := range errorEntries {
		fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
	}
}
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}

func parseLogLine(line string) (*LogEntry, error) {
    pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(line)

    if matches == nil {
        return nil, fmt.Errorf("invalid log format")
    }

    timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
    if err != nil {
        return nil, err
    }

    return &LogEntry{
        Timestamp: timestamp,
        Level:     matches[2],
        Message:   matches[3],
    }, nil
}

func filterLogsByLevel(entries []LogEntry, level string) []LogEntry {
    var filtered []LogEntry
    for _, entry := range entries {
        if strings.EqualFold(entry.Level, level) {
            filtered = append(filtered, entry)
        }
    }
    return filtered
}

func readLogFile(filename string) ([]LogEntry, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var entries []LogEntry
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err == nil {
            entries = append(entries, *entry)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return entries, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <logfile>")
        os.Exit(1)
    }

    entries, err := readLogFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error reading log file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Total log entries: %d\n", len(entries))

    errorLogs := filterLogsByLevel(entries, "ERROR")
    fmt.Printf("Error entries: %d\n", len(errorLogs))

    for _, entry := range errorLogs {
        fmt.Printf("[%s] %s\n", entry.Timestamp.Format("15:04:05"), entry.Message)
    }
}