
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func xorCipher(data []byte, key string) []byte {
	result := make([]byte, len(data))
	keyLen := len(key)
	for i := range data {
		result[i] = data[i] ^ key[i%keyLen]
	}
	return result
}

func processFile(inputPath, outputPath, key string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	inputInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	data := make([]byte, inputInfo.Size())
	_, err = io.ReadFull(inputFile, data)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	processedData := xorCipher(data, key)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(processedData)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func main() {
	inputFile := flag.String("input", "", "Path to input file")
	outputFile := flag.String("output", "", "Path to output file")
	key := flag.String("key", "defaultKey123", "Encryption key")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: file_encryptor -input <input> -output <output> [-key <key>]")
		os.Exit(1)
	}

	err := processFile(*inputFile, *outputFile, *key)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File processed successfully: %s -> %s\n", *inputFile, *outputFile)
}