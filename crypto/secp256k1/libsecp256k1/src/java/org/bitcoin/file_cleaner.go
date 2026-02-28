package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_cleaner.go <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := inputFile + ".deduped"

	seen := make(map[string]bool)

	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)

	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				fmt.Printf("Error writing to output file: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	writer.Flush()
	fmt.Printf("Duplicate lines removed. Output saved to: %s\n", outputFile)
}