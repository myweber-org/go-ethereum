package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption_tool.go <command>")
		fmt.Println("Commands: encrypt, decrypt, genkey")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "genkey":
		key, err := generateRandomKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

	case "encrypt":
		if len(os.Args) != 4 {
			fmt.Println("Usage: go run file_encryption_tool.go encrypt <input_file> <key_base64>")
			os.Exit(1)
		}

		inputFile := os.Args[2]
		keyBase64 := os.Args[3]

		key, err := base64.StdEncoding.DecodeString(keyBase64)
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			os.Exit(1)
		}

		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		encrypted, err := encryptData(data, key)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}

		outputFile := inputFile + ".enc"
		if err := os.WriteFile(outputFile, encrypted, 0644); err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Encrypted file saved as: %s\n", outputFile)

	case "decrypt":
		if len(os.Args) != 4 {
			fmt.Println("Usage: go run file_encryption_tool.go decrypt <encrypted_file> <key_base64>")
			os.Exit(1)
		}

		inputFile := os.Args[2]
		keyBase64 := os.Args[3]

		key, err := base64.StdEncoding.DecodeString(keyBase64)
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			os.Exit(1)
		}

		encryptedData, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		decrypted, err := decryptData(encryptedData, key)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}

		outputFile := "decrypted_" + inputFile[:len(inputFile)-4]
		if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Decrypted file saved as: %s\n", outputFile)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}