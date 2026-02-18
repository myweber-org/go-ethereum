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

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <logfile>")
        os.Exit(1)
    }
    
    file, err := os.Open(os.Args[1])
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()
    
    var entries []LogEntry
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            continue
        }
        entries = append(entries, *entry)
    }
    
    if len(os.Args) > 2 {
        levelFilter := os.Args[2]
        entries = filterLogsByLevel(entries, levelFilter)
    }
    
    for _, entry := range entries {
        fmt.Printf("%s [%s] %s\n", entry.Timestamp, entry.Level, entry.Message)
    }
}