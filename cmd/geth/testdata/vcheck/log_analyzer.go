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
	Source    string
}

type LogStats struct {
	TotalEntries int
	ErrorCount   int
	WarnCount    int
	InfoCount    int
	Sources      map[string]int
}

func parseLogLine(line string) (LogEntry, error) {
	pattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (\w+): (.+)$`)
	matches := pattern.FindStringSubmatch(line)

	if len(matches) != 5 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     matches[2],
		Message:   matches[4],
		Source:    matches[3],
	}, nil
}

func analyzeLogs(filePath string) (LogStats, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return LogStats{}, err
	}
	defer file.Close()

	stats := LogStats{
		Sources: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		stats.TotalEntries++
		stats.Sources[entry.Source]++

		switch strings.ToUpper(entry.Level) {
		case "ERROR":
			stats.ErrorCount++
		case "WARN":
			stats.WarnCount++
		case "INFO":
			stats.InfoCount++
		}
	}

	return stats, scanner.Err()
}

func printReport(stats LogStats) {
	fmt.Println("=== Log Analysis Report ===")
	fmt.Printf("Total entries: %d\n", stats.TotalEntries)
	fmt.Printf("Errors: %d\n", stats.ErrorCount)
	fmt.Printf("Warnings: %d\n", stats.WarnCount)
	fmt.Printf("Info messages: %d\n", stats.InfoCount)
	fmt.Println("\nSources breakdown:")
	for source, count := range stats.Sources {
		fmt.Printf("  %s: %d entries\n", source, count)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	stats, err := analyzeLogs(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}

	printReport(stats)
}