
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %w", err)
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %w", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(salt); err != nil {
		return fmt.Errorf("salt write failed: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return fmt.Errorf("iv generation failed: %w", err)
	}

	if _, err := outputFile.Write(iv); err != nil {
		return fmt.Errorf("iv write failed: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %w", err)
	}
	defer inputFile.Close()

	salt := make([]byte, 16)
	if _, err := io.ReadFull(inputFile, salt); err != nil {
		return fmt.Errorf("salt read failed: %w", err)
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, iv); err != nil {
		return fmt.Errorf("iv read failed: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %w", err)
	}
	defer outputFile.Close()

	stream := cipher.NewCTR(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inputFile}

	if _, err := io.Copy(outputFile, reader); err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	var err error
	switch mode {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, passphrase)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, passphrase)
	default:
		err = errors.New("invalid mode, use 'encrypt' or 'decrypt'")
	}

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s operation completed successfully\n", mode)
}