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

type LogSummary struct {
	TotalEntries int
	LevelCounts  map[string]int
	Errors       []string
}

func parseLogLine(line string) (LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     matches[2],
		Message:   matches[3],
	}, nil
}

func analyzeLogFile(filename string) (LogSummary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return LogSummary{}, err
	}
	defer file.Close()

	summary := LogSummary{
		LevelCounts: make(map[string]int),
		Errors:      []string{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		summary.TotalEntries++
		summary.LevelCounts[entry.Level]++

		if entry.Level == "ERROR" {
			summary.Errors = append(summary.Errors, entry.Message)
		}
	}

	return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
	fmt.Printf("Log Analysis Summary:\n")
	fmt.Printf("Total entries: %d\n", summary.TotalEntries)
	fmt.Printf("\nLevel distribution:\n")
	for level, count := range summary.LevelCounts {
		percentage := float64(count) / float64(summary.TotalEntries) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", level, count, percentage)
	}

	if len(summary.Errors) > 0 {
		fmt.Printf("\nRecent errors (%d found):\n", len(summary.Errors))
		for i, err := range summary.Errors {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", len(summary.Errors)-5)
				break
			}
			fmt.Printf("  - %s\n", strings.TrimSpace(err))
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	summary, err := analyzeLogFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing log file: %v\n", err)
		os.Exit(1)
	}

	printSummary(summary)
}