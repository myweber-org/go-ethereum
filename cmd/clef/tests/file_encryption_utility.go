package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	saltSize       = 16
	nonceSize      = 12
	keyIterations  = 100000
	keyLength      = 32
)

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	for i := 0; i < keyIterations; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)[:keyLength]
}

func encryptData(plaintext []byte, password string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	combined := append(salt, nonce...)
	combined = append(combined, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func decryptData(encoded string, password string) ([]byte, error) {
	combined, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if len(combined) < saltSize+nonceSize {
		return nil, errors.New("invalid encoded data")
	}

	salt := combined[:saltSize]
	nonce := combined[saltSize : saltSize+nonceSize]
	ciphertext := combined[saltSize+nonceSize:]

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func processFile(inputPath, outputPath, password string, encrypt bool) error {
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var result []byte
	var resultStr string
	if encrypt {
		resultStr, err = encryptData(inputData, password)
		if err != nil {
			return err
		}
		result = []byte(resultStr)
	} else {
		result, err = decryptData(string(inputData), password)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(outputPath, result, 0644)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file> <password>")
		os.Exit(1)
	}

	operation := strings.ToLower(os.Args[1])
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	password := os.Args[4]

	var encrypt bool
	switch operation {
	case "encrypt":
		encrypt = true
	case "decrypt":
		encrypt = false
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err := processFile(inputFile, outputFile, password, encrypt); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully. Output saved to: %s\n", outputFile)
}