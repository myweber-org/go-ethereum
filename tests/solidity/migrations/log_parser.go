
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

func parseLogLine(line string) *LogEntry {
    re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)`)
    matches := re.FindStringSubmatch(line)
    
    if matches == nil {
        return nil
    }
    
    return &LogEntry{
        Timestamp: matches[1],
        Level:     matches[2],
        Message:   matches[3],
    }
}

func filterErrors(entries []LogEntry) []LogEntry {
    var errors []LogEntry
    for _, entry := range entries {
        if strings.ToUpper(entry.Level) == "ERROR" {
            errors = append(errors, entry)
        }
    }
    return errors
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <logfile>")
        os.Exit(1)
    }
    
    filename := os.Args[1]
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()
    
    var entries []LogEntry
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        entry := parseLogLine(scanner.Text())
        if entry != nil {
            entries = append(entries, *entry)
        }
    }
    
    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }
    
    errorEntries := filterErrors(entries)
    
    fmt.Printf("Total log entries: %d\n", len(entries))
    fmt.Printf("Error entries: %d\n\n", len(errorEntries))
    
    for _, errEntry := range errorEntries {
        fmt.Printf("[%s] %s\n", errEntry.Timestamp, errEntry.Message)
    }
}