package main

import (
	"bufio"
	"fmt"
	"os"
)

func ProcessFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

func main() {
	lines, err := ProcessFileLines("input.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, line := range lines {
		fmt.Printf("Line %d: %s\n", i+1, line)
	}
}