package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_cleaner.go <input_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := inputFile + ".cleaned"

	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outFile.Close()

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)

	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			writer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	writer.Flush()
	fmt.Printf("Duplicate lines removed. Cleaned file saved as: %s\n", outputFile)
}