package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func xorCipher(input []byte, key string) []byte {
	keyBytes := []byte(key)
	output := make([]byte, len(input))
	for i := range input {
		output[i] = input[i] ^ keyBytes[i%len(keyBytes)]
	}
	return output
}

func processFile(inputPath, outputPath, key string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	inputData, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	outputData := xorCipher(inputData, key)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(outputData)
	return err
}

func main() {
	inputFile := flag.String("input", "", "Path to input file")
	outputFile := flag.String("output", "", "Path to output file")
	key := flag.String("key", "secret", "Encryption key")
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

	fmt.Println("File processed successfully")
}