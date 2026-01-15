package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
)

type LogParser struct {
    pattern *regexp.Regexp
}

func NewLogParser(regex string) (*LogParser, error) {
    compiled, err := regexp.Compile(regex)
    if err != nil {
        return nil, err
    }
    return &LogParser{pattern: compiled}, nil
}

func (lp *LogParser) ParseFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var matches []string
    scanner := bufio.NewScanner(file)
    lineNumber := 0

    for scanner.Scan() {
        lineNumber++
        line := scanner.Text()
        if lp.pattern.MatchString(line) {
            matches = append(matches, fmt.Sprintf("Line %d: %s", lineNumber, line))
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return matches, nil
}

func main() {
    parser, err := NewLogParser(`(?i)error|warning|fail`)
    if err != nil {
        fmt.Printf("Failed to create parser: %v\n", err)
        return
    }

    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <filename>")
        return
    }

    matches, err := parser.ParseFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error parsing file: %v\n", err)
        return
    }

    fmt.Printf("Found %d matching lines:\n", len(matches))
    for _, match := range matches {
        fmt.Println(match)
    }
}